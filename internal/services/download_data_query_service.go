package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
	"time"

	pb "github.com/DIMO-Network/shared/api/users"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const presignDurationHours = 24 * time.Hour

func (uds *UserDataService) UserDataJSONS3(user, key, start, end, ipfsAddress string, ipfs bool) error {
	query := uds.formatUserDataRequest(user, start, end)
	requested := time.Now().Format("2006-01-02 15:04:05")
	respSize := query.Size

	var ud UserData
	ud.User = user
	ud.RequestTimestamp = requested
	ud.RangeStart = start
	ud.RangeEnd = requested

	for respSize == query.Size {
		response, err := uds.executeESQuery(query)
		if err != nil {
			uds.log.Err(err).Msg("user data download: unable to query elasticsearch")
			return nil
		}

		respSize = int(gjson.Get(response, "hits.hits.#").Int())
		ud.DeviceID = gjson.Get(response, "hits.hits.0._source.data.device.device_id").String()

		data := make([]map[string]interface{}, respSize)
		err = json.Unmarshal([]byte(gjson.Get(response, "hits.hits").Raw), &data)
		if err != nil {
			uds.log.Err(err).Msg("user data download: unable to unmarshal data")
			return nil
		}

		ud.Data = append(ud.Data, data...)
		sA := gjson.Get(response, fmt.Sprintf("hits.hits.%d.sort.0", respSize-1))
		query.SearchAfter = []string{sA.String()}
	}

	keyName := "userDownloads/" + user + "/" + time.Now().Format(time.RFC3339) + ".json"
	s3link, err := uds.uploadUserData(ud, keyName)
	if err != nil {
		return nil
	}

	err = uds.sendEmail(user, s3link)
	if err != nil {
		return err
	}
	return nil
}

func (uds *UserDataService) formatUserDataRequest(user, rangestart, rangeend string) eSQuery {
	query := eSQuery{}
	query.Size = 10000
	query.formatESQuerySort(map[string]string{"data.timestamp": "desc"})
	query.formatESQueryFilterMust(map[string]string{"subject": user})
	query.formatESQueryFilterRange("data.timestamp", map[string]string{"gte": rangestart, "lte": rangeend})
	query.excludeFields([]string{"data.makeSlug", "data.modelSlug"})
	return query
}

type SearchUserData struct {
	Size        int     `json:"size"`
	Sort        sortBy  `json:"sort"`
	Filter      filter  `json:"query"`
	SearchAfter []int64 `json:"search_after,omitempty"`
}

type sortBy struct {
	DataTimestamp string `json:"data.timestamp"`
}

type UserData struct {
	User             string                   `json:"user"`
	RangeStart       string                   `json:"start"`
	RangeEnd         string                   `json:"end"`
	RequestTimestamp string                   `json:"requestTimestamp"`
	DeviceID         string                   `json:"deviceID"`
	Data             []map[string]interface{} `json:"data,omitempty"`
	EncryptedData    string                   `json:"encryptedData,omitempty"`
	IPFS             string                   `json:"ipfsAddress,omitempty"`
}

func (uds *UserDataService) sendEmail(user, downloadLink string) error {

	userEmail, err := uds.getVerifiedEmailAddress(user)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", uds.settings.EmailUsername, uds.settings.EmailPassword, uds.settings.EmailHost)
	addr := fmt.Sprintf("%s:%s", uds.settings.EmailHost, uds.settings.EmailPort)

	var partsBuffer bytes.Buffer
	w := multipart.NewWriter(&partsBuffer)
	defer w.Close() //nolint

	p, err := w.CreatePart(textproto.MIMEHeader{"Content-Type": {"text/plain"}, "Content-Transfer-Encoding": {"quoted-printable"}})
	if err != nil {
		return err
	}

	ptMessage := fmt.Sprintf("Hi,\r\n\r\nUse the following link(s) to download your requested data. These links will expire in 24 hours:\n\t%s\r\n\n", downloadLink)
	pw := quotedprintable.NewWriter(p)
	if _, err := pw.Write([]byte(ptMessage)); err != nil {
		return err
	}
	pw.Close()
	h, err := w.CreatePart(textproto.MIMEHeader{"Content-Type": {"text/html"}, "Content-Transfer-Encoding": {"quoted-printable"}})
	if err != nil {
		return err
	}

	hw := quotedprintable.NewWriter(h)
	htmlMessage := fmt.Sprintf(`<a href="%s">Click to download</a>`, downloadLink)
	if err := uds.emailTemplate.Execute(hw, template.HTML(htmlMessage)); err != nil {
		return err
	}
	hw.Close()
	var buffer bytes.Buffer
	buffer.WriteString("From: DIMO <" + uds.settings.EmailUsername + ">\r\n" +
		"To: " + userEmail + "\r\n" +
		"Subject: [DIMO] User Data Download\r\n" +
		"Content-Type: multipart/alternative; boundary=\"" + w.Boundary() + "\"\r\n" +
		"\r\n")
	if _, err := partsBuffer.WriteTo(&buffer); err != nil {
		return err
	}

	return smtp.SendMail(addr, auth, uds.settings.EmailFrom, []string{userEmail}, buffer.Bytes())
}

func (uds *UserDataService) getVerifiedEmailAddress(userID string) (string, error) {

	conn, err := grpc.Dial(uds.settings.UsersAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		uds.log.Err(err).Msg("failed to create users API client.")
		return "", nil
	}
	defer conn.Close()

	usersClient := pb.NewUserServiceClient(conn)
	user, err := usersClient.GetUser(context.Background(), &pb.GetUserRequest{Id: userID})
	if err != nil {
		return "", err
	}

	if user.EmailAddress == nil {
		uds.log.Error().Str("userId", user.Id).Msg("verified email address for user not found")
		emailNotFoundError := errors.New("verified email address for user not found")
		return "", emailNotFoundError
	}

	return *user.EmailAddress, nil
}

func (uds *UserDataService) putObjectS3(bucketname, keyname string, data []byte, svc *s3.S3) error {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(keyname),
		Body:   bytes.NewReader(data),
	}
	_, err := svc.PutObject(params)
	return err

}

func (uds *UserDataService) generatePreSignedURL(bucketname, keyName string, session *s3.S3, expiration time.Duration) (string, error) {
	req, _ := session.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(keyName),
	})
	return req.Presign(expiration)
}

func (uds *UserDataService) uploadUserData(ud UserData, keyName string) (string, error) {
	dataBytes, err := json.Marshal(ud)
	if err != nil {
		return "", err
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(uds.settings.AWSDefaultRegion),
		Credentials: credentials.NewStaticCredentials(uds.settings.AWSAccessKeyID, uds.settings.AWSSecretAccessKey, ""),
	})
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)
	err = uds.putObjectS3(uds.settings.AWSBucketName, keyName, dataBytes, svc)
	if err != nil {
		return "", err
	}
	return uds.generatePreSignedURL(uds.settings.AWSBucketName, keyName, svc, presignDurationHours)
}

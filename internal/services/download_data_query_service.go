package services

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/tidwall/gjson"
)

const presignDurationHours time.Duration = 24 * time.Hour

func (uds *UserDataService) DownloadUserData(user, key, start, end, ipfsAddress string, ipfs bool) (string, error) {

	query := uds.formatUserDataRequest(user, start, end)
	requested := time.Now().Format("2006-01-02 15:04:05")
	respSize := query.Size
	var dataDownloadLinks []string

	for respSize == query.Size {
		var ud UserData
		ud.User = user
		ud.RequestTimestamp = requested

		response, err := uds.executeESQuery(query)
		if err != nil {
			uds.log.Err(err).Msg("user data download: unable to query elasticsearch")
			return "", err
		}

		respSize = int(gjson.Get(response, "hits.hits.#").Int())

		ud.RangeStart = gjson.Get(response, "hits.hits.0.data.timestamp").Time().String()
		ud.RangeEnd = gjson.Get(response, fmt.Sprintf("hits.hits.%d.data.timestamp", respSize-1)).Time().String()

		ud.Data = make([]map[string]interface{}, respSize)
		err = json.Unmarshal([]byte(gjson.Get(response, "hits.hits").Raw), &ud.Data)
		if err != nil {
			uds.log.Err(err).Msg("user data download: unable to unmarshal data")
			return "", err
		}

		// comfirm what keyname and bucketname should be
		keyName := "userDownloads/" + user + "/" + time.Now().Format(time.RFC3339) + ".json"
		s3link, err := uds.uploadUserData(ud, keyName)
		if err != nil {
			return "", err
		}
		dataDownloadLinks = append(dataDownloadLinks, s3link)

		sA := gjson.Get(response, fmt.Sprintf("hits.hits.%d.sort.0", respSize-1))
		query.SearchAfter = []string{sA.String()}
	}

	err := uds.sendEmail(user, dataDownloadLinks)
	if err != nil {
		return "", err
	}
	// if key != "" {
	// 	var err error
	// 	bts, err := json.Marshal(ud.Data)
	// 	if err != nil {
	// 		aqs.log.Err(err).Msg("user data download: unable to marshal data")
	// 		return ud, err
	// 	}
	// 	ud.EncryptedData, err = encrypt(bts, key)
	// 	if err != nil {
	// 		return ud, err
	// 	}
	// 	ud.Data = nil
	// }

	// if ipfs {
	// 	url, err := uploadIPFS(ud.EncryptedData, ipfsAddress)
	// 	if err != nil {
	// 		return ud, err
	// 	}
	// 	ud.IPFS = url
	// 	ud.EncryptedData = ""
	// 	return ud, nil
	// }
	return "check user email address for download links", nil
}

func (uds *UserDataService) formatUserDataRequest(user, rangestart, rangeend string) eSQuery {
	var searchAfter string
	query := eSQuery{}
	query.Size = 10000
	query.formatESQuerySort(map[string]string{"data.timestamp": "desc"})
	query.formatESQueryFilterMust(map[string]string{"subject": user})
	query.formatESQueryFilterRange("data.timestamp", map[string]string{"gte": rangestart, "lte": rangeend})
	query.excludeFields([]string{"data.makeSlug", "data.modelSlug"})
	// set empty search after so this can be updated later
	query.SearchAfter = append(query.SearchAfter, searchAfter)
	return query
}

func uploadIPFS(encryptedData, ipfsAddress string) (string, error) {
	if encryptedData == "" {
		invalidUploadError := errors.New("failed to upload to ipfs: data must be encrypted")
		return "", invalidUploadError
	}

	sh := shell.NewShell(ipfsAddress)
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(encryptedData)
	if err != nil {
		return "", err
	}
	cid, err := sh.Add(&buf)
	if err != nil {
		return "", err
	}
	// data available at link
	url := fmt.Sprintf("https://ipfs.io/ipfs/%s", cid)

	return url, nil
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
	Data             []map[string]interface{} `json:"data,omitempty"`
	EncryptedData    string                   `json:"encryptedData,omitempty"`
	DecryptedData    []map[string]interface{} `json:"decryptedData,omitempty"`
	IPFS             string                   `json:"ipfsAddress,omitempty"`
}

func createHash(key string) string {
	// use a more secure encryption method here?
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) (string, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(ciphertext), nil
}

// func decrypt(data string, passphrase string) ([]byte, error) {

// 	dataBytes, err := hex.DecodeString(data)
// 	key := []byte(createHash(passphrase))
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return []byte{}, err
// 	}
// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return []byte{}, err
// 	}
// 	nonceSize := gcm.NonceSize()
// 	nonce, ciphertext := dataBytes[:nonceSize], dataBytes[nonceSize:]
// 	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
// 	if err != nil {
// 		return []byte{}, err
// 	}
// 	return plaintext, nil
// }

func (uds *UserDataService) sendEmail(user string, links []string) error {

	userEmail, err := getVerifiedEmailAddress(user)

	auth := smtp.PlainAuth("", uds.settings.EmailUsername, uds.settings.EmailPassword, uds.settings.EmailHost)
	addr := fmt.Sprintf("%s:%s", uds.settings.EmailHost, uds.settings.EmailPort)

	var partsBuffer bytes.Buffer
	w := multipart.NewWriter(&partsBuffer)
	defer w.Close() //nolint

	p, err := w.CreatePart(textproto.MIMEHeader{"Content-Type": {"text/plain"}, "Content-Transfer-Encoding": {"quoted-printable"}})
	if err != nil {
		return err
	}

	pw := quotedprintable.NewWriter(p)
	message := "Hi,\r\n\r\nUse the following link(s) to download your requested data. These links will expire in 24 hours: "
	for _, link := range links {
		message += link + "\r\n"
	}
	message += "\n\n"

	if _, err := pw.Write([]byte(message)); err != nil {
		return err
	}
	pw.Close()

	var buffer bytes.Buffer
	buffer.WriteString("From: DIMO <" + uds.settings.EmailUsername + ">\r\n" +
		"To: " + userEmail + "\r\n" +
		"Subject: [DIMO] User Download\r\n" +
		"Content-Type: text/plain charset=utf-8; boundary=\"" + w.Boundary() + "\"\r\n" +
		"\r\n")
	if _, err := partsBuffer.WriteTo(&buffer); err != nil {
		return err
	}

	return smtp.SendMail(addr, auth, uds.settings.EmailFrom, []string{userEmail}, buffer.Bytes())
}

func getVerifiedEmailAddress(user string) (string, error) {

	// is there a grpc endpoint that can return the user email?
	// otherwise grab user email from db

	// user, err := models.Users(
	// 	models.UserWhere.ID.EQ(userID),
	// 	qm.Load(models.UserRels.Referrals),
	// ).One(c.Context(), tx)
	// if err != nil {
	// 	if !errors.Is(err, sql.ErrNoRows) {
	// 		return nil, err
	// 	}
	// }
	return "", nil
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

	creds := credentials.NewStaticCredentials(uds.settings.AWSAccessKeyID, uds.settings.AWSSecretAccessKey, "")
	cfg := aws.NewConfig().WithRegion(uds.settings.AWSDefaultRegion).WithCredentials(creds)
	sess, err := session.NewSession(cfg)
	svc := s3.New(sess)

	err = uds.putObjectS3(uds.settings.AWSBucketName, keyName, dataBytes, svc)
	if err != nil {
		return "", err
	}
	return uds.generatePreSignedURL(uds.settings.AWSBucketName, keyName, svc, presignDurationHours)
}

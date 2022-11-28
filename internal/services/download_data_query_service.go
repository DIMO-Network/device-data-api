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
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/tidwall/gjson"
)

func (aqs *AggregateQueryService) DownloadUserData(user, key, start, end, ipfsAddress string, ipfs bool) (UserData, error) {
	var searchAfter string
	query := eSQuery{}
	query.Size = 1000
	query.formatESQuerySort(map[string]string{"data.timestamp": "desc"})
	query.formatESQueryFilterMust(map[string]string{"subject": user})
	query.formatESQueryFilterRange("data.timestamp", map[string]string{"gte": start, "lte": end})
	if searchAfter != "" {
		query.SearchAfter = append(query.SearchAfter, searchAfter)
	}
	var ud UserData
	ud.User = user
	ud.RangeStart = start
	ud.RangeEnd = end
	ud.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
	respSize := query.Size
	for respSize == query.Size {
		response, err := aqs.executeESQuery(query)
		if err != nil {
			aqs.log.Err(err).Msg("user data download: unable to query elasticsearch")
		}

		respSize = int(gjson.Get(response, "hits.hits.#").Int())
		data := make([]map[string]interface{}, respSize)
		err = json.Unmarshal([]byte(gjson.Get(response, "hits.hits").Raw), &data)
		if err != nil {
			aqs.log.Err(err).Msg("user data download: unable to unmarshal data")
		}
		ud.Data = append(ud.Data, data...)
		sA := gjson.Get(response, fmt.Sprintf("hits.hits.%d.sort.0", respSize-1))
		query.SearchAfter = []string{sA.String()}
	}

	if key != "" {
		var err error
		bts, err := json.Marshal(ud.Data)
		if err != nil {
			aqs.log.Err(err).Msg("user data download: unable to query elasticsearch")
		}
		ud.EncryptedData, err = encrypt(bts, key)
		if err != nil {
			return ud, err
		}
		ud.Data = nil
	}

	if ipfs {
		url, err := uploadIPFS(ud.EncryptedData, ipfsAddress)
		if err != nil {
			return ud, err
		}
		ud.IPFS = url
		ud.EncryptedData = ""
		return ud, nil
	}
	return ud, nil
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
	User          string                   `json:"user"`
	RangeStart    string                   `json:"start"`
	RangeEnd      string                   `json:"end"`
	TimeStamp     string                   `json:"timestamp"`
	Data          []map[string]interface{} `json:"data,omitempty"`
	EncryptedData string                   `json:"encryptedData,omitempty"`
	DecryptedData []map[string]interface{} `json:"decryptedData,omitempty"`
	IPFS          string                   `json:"ipfsAddress,omitempty"`
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

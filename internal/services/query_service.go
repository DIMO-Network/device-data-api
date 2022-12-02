package services

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

const presignDuration = 24

type UserDataService struct {
	es        *elasticsearch.TypedClient
	indexName string
	log       *zerolog.Logger
}

type UserData struct {
	User             string                   `json:"user"`
	RequestTimestamp string                   `json:"requestTimestamp"`
	DeviceID         string                   `json:"deviceID"`
	Data             []map[string]interface{} `json:"data,omitempty"`
}

func NewAggregateQueryService(es *elasticsearch.TypedClient, settings *config.Settings, log *zerolog.Logger) *UserDataService {
	return &UserDataService{es: es, indexName: settings.DeviceDataIndexName, log: log}
}

func (uds *UserDataService) executeESQuery(q *search.Request) (string, error) {
	res, err := uds.es.Search().
		Index(uds.indexName).
		Request(q).
		Do(context.TODO())
	if err != nil {
		uds.log.Err(err).Msg("Could not query Elasticsearch")
		return "", err
	}
	defer res.Body.Close()

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		uds.log.Err(err).Msg("Could not parse Elasticsearch response body")
		return "", err
	}
	response := string(responseBytes)

	if res.StatusCode != 200 {
		uds.log.Info().RawJSON("elasticsearchResponseBody", responseBytes).Msg("Error from Elastic.")

		err := fmt.Errorf("invalid status code when querying elastic: %d", res.StatusCode)
		return response, err
	}

	return response, nil
}

func (uds *UserDataService) UserDataJSONS3(userDeviceID string) (UserData, error) {
	query := uds.formatUserDataRequest(userDeviceID)
	requested := time.Now().Format(time.RFC3339)
	respSize := pageSize

	ud := UserData{
		User:             userDeviceID,
		RequestTimestamp: requested,
	}

	for respSize == pageSize {
		response, err := uds.executeESQuery(query)
		if err != nil {
			uds.log.Err(err).Msg("user data download: unable to query elasticsearch")
			return UserData{}, err
		}

		respSize = int(gjson.Get(response, "hits.hits.#").Int())
		ud.DeviceID = gjson.Get(response, "hits.hits.0._source.data.device.device_id").String()

		data := make([]map[string]interface{}, respSize)
		err = json.Unmarshal([]byte(gjson.Get(response, "hits.hits").Raw), &data)
		if err != nil {
			uds.log.Err(err).Msg("user data download: unable to unmarshal data")
			return UserData{}, err
		}

		ud.Data = append(ud.Data, data...)
		sA := gjson.Get(response, fmt.Sprintf("hits.hits.%d.sort.0", respSize-1))
		query.SearchAfter = []types.FieldValue{sA.String()}
	}

	return ud, nil
}

// Elastic maximum.
var pageSize = 10000

func (uds *UserDataService) formatUserDataRequest(userDeviceID string) *search.Request {
	query := &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{Match: map[string]types.MatchQuery{"subject": {Query: userDeviceID}}},
				},
			},
		},
		Sort: []types.SortCombinations{"data.timestamp"},
		Size: &pageSize,
	}
	return query
}

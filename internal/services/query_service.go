package services

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"reflect"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/customerio/go-customerio/v3"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/rs/zerolog"
)

//go:embed data_download_email_template.html
var rawDataDownloadEmail string

type UserDataService struct {
	settings      *config.Settings
	es            *elasticsearch.Client
	log           *zerolog.Logger
	emailTemplate *template.Template
	cioClient     *customerio.CustomerIO
}

func NewAggregateQueryService(es *elasticsearch.Client, log *zerolog.Logger, settings *config.Settings) *UserDataService {
	t := template.Must(template.New("data_download_email_template").Parse(rawDataDownloadEmail))
	var cioClient *customerio.CustomerIO
	return &UserDataService{es: es, log: log, settings: settings, emailTemplate: t, cioClient: cioClient}
}

func (uds *UserDataService) executeESQuery(q interface{}) (string, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(q); err != nil {
		return "", err
	}

	res, err := uds.es.Search(
		uds.es.Search.WithContext(context.Background()),
		uds.es.Search.WithIndex(uds.settings.ElasticIndex),
		uds.es.Search.WithBody(&buf),
	)
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
		uds.log.Info().RawJSON("elasticRequest", buf.Bytes()).Msg("request sent to elastics")

		err := fmt.Errorf("invalid status code when querying elastic: %d", res.StatusCode)
		return response, err
	}

	return response, nil
}

func (q *eSQuery) formatESQuerySort(sort map[string]string) {
	for fieldname, sortDirection := range sort {
		matchOn := term{}
		matchOn.Term = sortDirection
		termFilterType := reflect.TypeOf(matchOn)
		termFilterValue := reflect.ValueOf(matchOn)

		fs := []reflect.StructField{}
		for i := 0; i < termFilterType.NumField(); i++ {
			fs = append(fs, termFilterType.Field(i))
		}
		fs[0].Tag = reflect.StructTag(fmt.Sprintf(`json:"%s"`, fieldname))
		st2 := reflect.StructOf(fs)
		v2 := termFilterValue.Convert(st2)
		q.Sort = v2.Interface()
	}
}

func (q *eSQuery) formatESQueryFilterMust(must map[string]string) {
	for fieldname, matchterm := range must {
		matchOn := term{}
		matchOn.Term = matchterm
		termFilterType := reflect.TypeOf(matchOn)
		termFilterValue := reflect.ValueOf(matchOn)

		fs := []reflect.StructField{}
		for i := 0; i < termFilterType.NumField(); i++ {
			fs = append(fs, termFilterType.Field(i))
		}
		fs[0].Tag = reflect.StructTag(fmt.Sprintf(`json:"%s"`, fieldname))
		st2 := reflect.StructOf(fs)
		v2 := termFilterValue.Convert(st2)
		must := matchFilter{}
		must.Match = v2.Interface()
		q.Filter.Bool.Must = append(q.Filter.Bool.Must, must)
	}
}

func (q *eSQuery) formatESQueryFilterRange(rangefield string, rangeMap map[string]string) {
	tRange := dateRangeFormat{}
	if rangeMap["gte"] != "" {
		tRange.DataTimestamp.Gte = rangeMap["gte"]
	}
	if rangeMap["gt"] != "" {
		tRange.DataTimestamp.Gt = rangeMap["gt"]
	}
	if rangeMap["lte"] != "" {
		tRange.DataTimestamp.Lte = rangeMap["lte"]
	}
	if rangeMap["lt"] != "" {
		tRange.DataTimestamp.Lt = rangeMap["lt"]
	}

	tRangeType := reflect.TypeOf(tRange)
	tRangeValue := reflect.ValueOf(tRange)

	fs := []reflect.StructField{}
	for i := 0; i < tRangeType.NumField(); i++ {
		fs = append(fs, tRangeType.Field(i))
	}
	fs[0].Tag = reflect.StructTag(fmt.Sprintf(`json:"%s"`, rangefield))
	st2 := reflect.StructOf(fs)
	v2 := tRangeValue.Convert(st2)
	r := dateRange{}
	r.DataTimestamp = v2.Interface()
	q.Filter.Bool.Filter = append(q.Filter.Bool.Filter, r)

}

func (q *eSQuery) excludeFields(terms []string) {
	q.ResponseFields.Exclude = terms
}

type eSQuery struct {
	ResponseFields responseFields `json:"_source,omitempty"`
	Aggs           interface{}    `json:"aggs,omitempty"`
	Size           int            `json:"size"`
	Filter         filter         `json:"query,omitempty"`
	Fields         []interface{}  `json:"fields,omitempty"`
	Sort           interface{}    `json:"sort,omitempty"`
	SearchAfter    []string       `json:"search_after,omitempty"`
}

type filter struct {
	Match interface{} `json:"match,omitempty"`
	Bool  struct {
		Must   []interface{} `json:"must,omitempty"`
		Filter []interface {
		} `json:"filter,omitempty"`
	} `json:"bool,omitempty"`
}

type term struct {
	Term string `json:"term"`
}

type matchFilter struct {
	Match interface{} `json:"match"`
}

type dateRangeFormat struct {
	DataTimestamp struct {
		Format string      `json:"format,omitempty"`
		Gte    interface{} `json:"gte,omitempty"`
		Lte    interface{} `json:"lte,omitempty"`
		Gt     interface{} `json:"gt,omitempty"`
		Lt     interface{} `json:"lt,omitempty"`
	} `json:"data.timestamp"`
}

type dateRange struct {
	DataTimestamp interface{} `json:"range"`
}

type responseFields struct {
	Exclude []string `json:"exclude,omitempty"`
	Include []string `json:"include,omitempty"`
}

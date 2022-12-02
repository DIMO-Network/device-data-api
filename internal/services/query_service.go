package services

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"io"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/customerio/go-customerio/v3"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/rs/zerolog"
)

//go:embed data_download_email_template.html
var rawDataDownloadEmail string

type UserDataService struct {
	settings      *config.Settings
	es            *elasticsearch.TypedClient
	log           *zerolog.Logger
	emailTemplate *template.Template
	cioClient     *customerio.CustomerIO
}

func NewAggregateQueryService(es *elasticsearch.TypedClient, log *zerolog.Logger, settings *config.Settings) *UserDataService {
	t := template.Must(template.New("data_download_email_template").Parse(rawDataDownloadEmail))
	var cioClient *customerio.CustomerIO
	return &UserDataService{es: es, log: log, settings: settings, emailTemplate: t, cioClient: cioClient}
}

func (uds *UserDataService) executeESQuery(q *search.Request) (string, error) {
	res, err := uds.es.Search().
		Index(uds.settings.DeviceDataIndexName).
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

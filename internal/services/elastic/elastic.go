// Package elastic is responsible for handling queries to the backing elastic database.
package elastic

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/some"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/calendarinterval"
	"github.com/rs/zerolog"

	"github.com/DIMO-Network/device-data-api/internal/config"
)

const (
	defaultRetries = 5
)

// Service is a service for performing queries on elastic.
type Service struct {
	settings *config.Settings
	log      *zerolog.Logger
	esClient *elasticsearch.TypedClient
}

// New creates a newly configured elastic service.
func New(settings *config.Settings, logger *zerolog.Logger, caCert []byte) (*Service, error) {
	esConfig := elasticsearch.Config{
		Addresses:  []string{settings.ElasticSearchAnalyticsHost},
		Username:   settings.ElasticSearchAnalyticsUsername,
		Password:   settings.ElasticSearchAnalyticsPassword,
		CACert:     caCert,
		MaxRetries: defaultRetries,
	}
	es8Client, err := elasticsearch.NewTypedClient(esConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create elastic client: %w", err)
	}
	eClient := Service{
		esClient: es8Client,
		log:      logger,
		settings: settings,
	}
	return &eClient, nil
}

// ESClient returns the underlying elastic client used by the service.
//
// This function exists to maintain previous behavior for endpoint who are not yet using the elastic service.
func (s *Service) ESClient() *elasticsearch.TypedClient {
	return s.esClient
}

func (s *Service) GetTotalDailyDistanceDriven(ctx context.Context, tz, deviceID string) ([]byte, error) {
	query := &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{Match: map[string]types.MatchQuery{"subject": {Query: deviceID}}},
					{Range: map[string]types.RangeQuery{"data.timestamp": types.DateRangeQuery{Gte: some.String("now-13d/d"), TimeZone: &tz}}},
				},
			},
		},
		Size: some.Int(0),
		Aggregations: map[string]types.Aggregations{
			"days": {
				DateHistogram: &types.DateHistogramAggregation{
					Field:            some.String("data.timestamp"),
					CalendarInterval: &calendarinterval.Day,
					TimeZone:         &tz,
				},
				Aggregations: map[string]types.Aggregations{
					"min_odom": {
						Min: &types.MinAggregation{
							Field: some.String("data.odometer"),
						},
					},
					"max_odom": {
						Max: &types.MaxAggregation{
							Field: some.String("data.odometer"),
						},
					},
					// Code generation for buckets_path is broken as of 8.5.0
					// See https://github.com/elastic/go-elasticsearch/issues/570
				},
			},
		},
	}

	resp, err := s.esClient.Search().Index(s.settings.DeviceDataIndexName).Request(query).Perform(ctx)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() // nolint

	if c := resp.StatusCode; c != 200 {
		return nil, fmt.Errorf("got status code %d from Elastic message", c)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body %w", err)
	}

	return body, nil
}

func (s *Service) GetTotalDistanceDriven(ctx context.Context, deviceID string) ([]byte, error) {
	query := &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{Term: map[string]types.TermQuery{"subject": {Value: deviceID}}},
				},
			},
		},
		Size: some.Int(0),
		Aggregations: map[string]types.Aggregations{
			"max_odometer": {
				Max: &types.MaxAggregation{
					Field: some.String("data.odometer"),
				},
			},
			"min_odometer": {
				Min: &types.MinAggregation{
					Field: some.String("data.odometer"),
				},
			},
		},
	}

	res, err := s.esClient.Search().Index(s.settings.DeviceDataIndexName).Request(query).Perform(ctx)
	if err != nil {
		return nil, fmt.Errorf("error querying odometer %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status code %d from Elastic message", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// Package elastic is responsible for handling queries to the backing elastic database.
package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/shared/privileges"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

const (
	docByInterval  = "documents_by_interval"
	singleDoc      = "select_single_doc"
	timeField      = "time"
	subjectField   = "subject"
	defaultRetries = 5
)

// ErrInvalidParams is returned when the parameters for the history query are invalid.
var ErrInvalidParams = fmt.Errorf("invalid parameters")

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

// GetHistoryParams defines the parameters for retrieving history from elastic.
type GetHistoryParams struct {
	// StartTime is the start time of the history.
	StartTime time.Time
	// EndTime is the end time of the history.
	EndTime time.Time
	// DeviceID is the ID of the device.
	DeviceID string
	// PrivilegeIDs is the list of privilege IDs that the user has.
	PrivilegeIDs []privileges.Privilege
	// Buckets is the number of time intervals to divide the history into.
	Buckets int
}

// GetHistory retrieves the history of a device from elastic. The history is divided into buckets and one data point is selected from each bucket.
// The data points are selected by the first data point in the bucket that has a vehicle field.
// The result is a list of data points.
func (s *Service) GetHistory(ctx context.Context, params GetHistoryParams) ([]json.RawMessage, error) {
	if params.Buckets < 1 {
		return nil, fmt.Errorf("non-positive number of buckets: %w", ErrInvalidParams)
	}
	// build the query
	req := buildHistoryQuery(params)

	// perform the request
	res, err := s.esClient.Search().Index(s.settings.DeviceDataIndexNameV2).Request(&req).Perform(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to complete request: %w", err)
	}

	//nolint:errcheck // we don't care about the error closing the body.
	defer res.Body.Close()

	// check for errors any status code over 299 is an error
	if res.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(res.Body) // we ignore the error here because we are already in an error state
		return nil, fmt.Errorf("got status code %d from Elastic message = %s", res.StatusCode, string(body))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// depending on the start and end time we may get len(results) <=  buckets + 1
	// pass the max number of values to getHitsFromHistoryResponse to ensure we only get the desired number of values
	return getHitsFromHistoryResponse(body, params.Buckets), nil
}

// SetDefaultHistoryParams sets the default values for the history params.
// If the buckets are less than or equal to 0, it will be set to 1000.
// If the end time is zero, it will be set to the current time.
// If the start time is zero, it will be set to 2 weeks from the end time.
func (g *GetHistoryParams) SetDefaultHistoryParams() {
	if g.Buckets <= 0 {
		g.Buckets = 1000
	}
	if g.EndTime.IsZero() {
		g.EndTime = time.Now()
	}

	if g.StartTime.IsZero() {
		g.StartTime = g.EndTime.Add(-time.Hour * 24 * 14) // default to 2 weeks ago
	}
}

// buildHistoryQuery builds the query to retrieve the history of a device from elastic.
func buildHistoryQuery(params GetHistoryParams) search.Request {
	// convert the start and end times to milliseconds since epoch
	startArg := strconv.Itoa(int(params.StartTime.UnixMilli()))
	endArg := strconv.Itoa(int(params.EndTime.UnixMilli()))

	// calculate the interval in milliseconds
	interval := calculateIntervalMS(params.StartTime, params.EndTime, params.Buckets)
	intervalArg := strconv.Itoa(interval) + "ms"

	min := 1
	filter := getHistoricalPrivilegeFilter(params.PrivilegeIDs)
	return search.Request{
		// get documents that match the device ID, are within the time range, and have a vehicle field
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{Term: map[string]types.TermQuery{subjectField: {Value: params.DeviceID}}},
					{Range: map[string]types.RangeQuery{timeField: types.DateRangeQuery{Gte: &startArg, Lte: &endArg}}},
					{Exists: &types.ExistsQuery{Field: "data.vehicle"}},
				},
			},
		},
		Aggregations: map[string]types.Aggregations{
			// docByInterval histogram aggregates the data into buckets based on the interval
			docByInterval: {
				DateHistogram: &types.DateHistogramAggregation{
					Field:         ptr(timeField),
					FixedInterval: &intervalArg,
					MinDocCount:   &min,
				},
				Aggregations: map[string]types.Aggregations{
					// singleDoc aggregation selects the first data point in each bucket
					// that matches the desired filter
					singleDoc: {
						TopHits: &types.TopHitsAggregation{
							Size:    ptr(1),
							Source_: &filter,
						},
					},
				},
			},
		},
		Size: ptr(0), // only return the aggregations
	}
}

// calculateIntervalMS calculates the desired interval in milliseconds. This interval will provide n buckets between the start and end times.
// if the interval is 0 or negative it will be clamped to 1ms.
func calculateIntervalMS(start, end time.Time, buckets int) int {
	duration := int(end.Sub(start).Milliseconds())
	interval := duration / buckets
	if interval < 1 {
		return 1 // 0 or negative interval is not supported, clamp to 1
	}
	return interval
}

// getHitsFromHistoryResponse extracts the hits from the response body.
// The maxValues parameter is used to limit the number of values returned.
// Due to the mapping function used by elaticsearch `bucket_key = Math.floor(time / interval) * interval
// the number of buckets returned may be 1 more than the number of buckets requested.
func getHitsFromHistoryResponse(body []byte, maxValues int) []json.RawMessage {
	retData := []json.RawMessage{}
	// This key gets the first data point in each bucket which should only contain 1 value
	fullKey := fmt.Sprintf("aggregations.%s.buckets.#.%s.hits.hits.0._source", docByInterval, singleDoc)
	gjson.GetBytes(body, fullKey).ForEach(func(key, value gjson.Result) bool {
		// if we have reached the max number of values, stop
		if len(retData) == maxValues {
			return false
		}
		retData = append(retData, resultToBytes(body, value))
		return true
	})
	return retData
}

// resultToBytes converts the result to a byte slice.
// this logic is the recommended way to extract the raw bytes from a gjson result.
// more info: https://github.com/tidwall/gjson/blob/6ee9f877d683381343bc998c137339c7ae908b86/README.md#working-with-bytes
func resultToBytes(body []byte, result gjson.Result) []byte {
	var raw []byte
	if result.Index > 0 {
		raw = body[result.Index : result.Index+len(result.Raw)]
	} else {
		raw = []byte(result.Raw)
	}
	return raw
}

// getHistoricalPrivilegeFilter returns the location filter based on the privileges.Vehicle
func getHistoricalPrivilegeFilter(privs []privileges.Privilege) types.SourceFilter {
	var filter types.SourceFilter

	canSeeLocation := slices.Contains(privs, privileges.VehicleAllTimeLocation)
	canSeeNonLocationData := slices.Contains(privs, privileges.VehicleNonLocationData)

	switch {
	case canSeeLocation && canSeeNonLocationData:
		// user can see everything
		// do nothing to the filter
	case canSeeLocation:
		// user can only see location data
		filter.Includes = append(filter.Includes, "*.misc.cell", "*.vehicle.currentLocation")
	case canSeeNonLocationData:
		// user can see everything except location data
		filter.Excludes = append(filter.Excludes, "*.misc.cell", "*.vehicle.currentLocation")
	default:
		// user can't see anything
		filter.Excludes = append(filter.Excludes, "*")
	}

	return filter
}

// ptr returns a pointer to the value passed in.
func ptr[T any](value T) *T {
	return &value
}

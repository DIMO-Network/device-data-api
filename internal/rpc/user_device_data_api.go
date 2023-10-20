package rpc

import (
	"context"

	"database/sql"
	"fmt"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/pkg/errors"
	smartcar "github.com/smartcar/go-sdk"

	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	internalmodel "github.com/DIMO-Network/device-data-api/internal/models"
	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/DIMO-Network/device-data-api/models"
	pb "github.com/DIMO-Network/device-data-api/pkg/grpc"
	"github.com/DIMO-Network/shared/db"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewUserDeviceData(dbs func() *db.ReaderWriter, logger *zerolog.Logger, deviceDefSvc services.DeviceDefinitionsAPIService, deviceStatusSvc services.DeviceStatusService) pb.UserDeviceDataServiceServer {
	return &userDeviceData{dbs: dbs, logger: logger, deviceDefSvc: deviceDefSvc, deviceStatusSvc: deviceStatusSvc}
}

type userDeviceData struct {
	pb.UserDeviceDataServiceServer
	dbs             func() *db.ReaderWriter
	logger          *zerolog.Logger
	deviceDefSvc    services.DeviceDefinitionsAPIService
	deviceStatusSvc services.DeviceStatusService
}

// ! TODO: need test for this
func (s *userDeviceData) GetUserDeviceData(ctx context.Context, req *pb.UserDeviceDataRequest) (*pb.UserDeviceDataResponse, error) {
	if isEmpty(req.UserDeviceId) || isEmpty(req.DeviceDefinitionId) {
		return nil, status.Error(codes.InvalidArgument, "UserDeviceId and DeviceDefinitionId are required")
	}
	deviceData, err := models.UserDeviceData(
		models.UserDeviceDatumWhere.UserDeviceID.EQ(req.UserDeviceId),
		models.UserDeviceDatumWhere.Signals.IsNotNull(),
		models.UserDeviceDatumWhere.UpdatedAt.GT(time.Now().Add(-14*24*time.Hour)),
	).All(ctx, s.dbs().Reader)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	if len(deviceData) == 0 {
		return nil, status.Error(codes.NotFound, "No status updates yet.")
	}
	var deviceStyleID *string
	if len(req.DeviceStyleId) > 0 {
		deviceStyleID = &req.DeviceStyleId
	} else {
		deviceStyleID = nil
	}

	ds := s.deviceStatusSvc.PrepareDeviceStatusInformation(ctx, deviceData, req.DeviceDefinitionId,
		deviceStyleID, req.PrivilegeIds) // up to caller to pass in correct privileges

	return &pb.UserDeviceDataResponse{
		Charging:             ds.Charging,
		FuelPercentRemaining: ds.FuelPercentRemaining,
		BatteryCapacity:      ds.BatteryCapacity,
		OilLevel:             ds.OilLevel,
		Odometer:             ds.Odometer,
		Latitude:             ds.Latitude,
		Longitude:            ds.Longitude,
		Range:                ds.Range,
		StateOfCharge:        ds.StateOfCharge,
		ChargeLimit:          ds.ChargeLimit,
		RecordUpdatedAt:      convertToTimestamp(ds.RecordUpdatedAt),
		RecordCreatedAt:      convertToTimestamp(ds.RecordCreatedAt),
		TirePressure:         convertTirePressure(ds.TirePressure),
		BatteryVoltage:       ds.BatteryVoltage,
		AmbientTemp:          ds.AmbientTemp,
	}, nil
}

func (s *userDeviceData) GetRawDeviceData(ctx context.Context, req *pb.RawDeviceDataRequest) (*pb.RawDeviceDataResponse, error) {
	if isEmpty(req.UserDeviceId) {
		return nil, status.Error(codes.InvalidArgument, "UserDeviceId is required")
	}

	query := make([]qm.QueryMod, 0)

	query = append(query, models.UserDeviceDatumWhere.UserDeviceID.EQ(req.UserDeviceId))

	if req.IntegrationId != nil && !isEmpty(*req.IntegrationId) {
		query = append(query, models.UserDeviceDatumWhere.IntegrationID.EQ(*req.IntegrationId))
	}

	deviceData, err := models.UserDeviceData(
		query...,
	).All(ctx, s.dbs().Reader)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	responseData := make([]*pb.RawDeviceDataResponseItem, len(deviceData))

	for i, d := range deviceData {
		responseData[i] = &pb.RawDeviceDataResponseItem{
			UserDeviceId:        d.UserDeviceID,
			IntegrationId:       d.IntegrationID,
			SignalsJsonData:     d.Signals.JSON,
			ErrorJsonData:       d.ErrorData.JSON,
			RecordUpdatedAt:     convertToTimestamp(&d.UpdatedAt),
			RecordCreatedAt:     convertToTimestamp(&d.CreatedAt),
			LastLocationEventAt: convertToTimestamp(&d.LastLocationEventAt.Time),
			LastObd2EventAt:     convertToTimestamp(&d.LastOdb2EventAt.Time),
		}
	}

	return &pb.RawDeviceDataResponse{
		Items: responseData,
	}, nil

}

func convertToTimestamp(goTime *time.Time) *timestamppb.Timestamp {
	if goTime == nil {
		return nil
	}
	timestamp := timestamppb.New(*goTime)
	return timestamp
}

func convertTirePressure(tp *smartcar.TirePressure) *pb.TirePressureResponse {
	if tp == nil {
		return nil
	}
	return &pb.TirePressureResponse{
		FrontLeft:  tp.FrontLeft,
		FrontRight: tp.FrontRight,
		BackLeft:   tp.BackLeft,
		BackRight:  tp.BackRight,
		DataAge:    tp.DataAge,
		RequestId:  tp.RequestID,
		UnitSystem: string(tp.UnitSystem),
	}
}

func (s *userDeviceData) GetSignals(ctx context.Context, req *pb.SignalRequest) (*pb.SignalResponse, error) {

	keys := make(map[string]string)

	keys["0"] = "property_id"
	keys["1"] = "device_make"
	keys["2"] = "model"
	keys["3"] = "year"

	queryMods := []qm.QueryMod{
		qm.Select(fmt.Sprintf("%s as name", keys[*req.Level]), "SUM(count) as total_count"),
		qm.GroupBy(keys[*req.Level]),
	}

	queryMods = append(queryMods, models.ReportVehicleSignalsEventsTrackingWhere.IntegrationID.EQ(req.IntegrationId))
	queryMods = append(queryMods, models.ReportVehicleSignalsEventsTrackingWhere.DateID.EQ(req.DateId))

	if *req.Level == "1" {
		if req.PropertyId == nil || *req.PropertyId == "" {
			return nil, status.Error(codes.InvalidArgument, "Invalid argument. PropertyId is required.")
		}
	}

	if *req.Level == "2" {
		if req.Make == nil || *req.Make == "" {
			return nil, status.Error(codes.InvalidArgument, "Invalid argument. Make is required.")
		}
	}

	if *req.Level == "3" {
		if req.Model == nil || *req.Model == "" {
			return nil, status.Error(codes.InvalidArgument, "Invalid argument. Model is required.")
		}
	}

	if *req.Level == "4" {
		if req.Year == nil || *req.Year == 0 {
			return nil, status.Error(codes.InvalidArgument, "Invalid argument. Year is required.")
		}
	}

	if req.PropertyId != nil && *req.PropertyId != "" {
		queryMods = append(queryMods, models.ReportVehicleSignalsEventsTrackingWhere.PropertyID.EQ(*req.PropertyId))
	}

	if req.Make != nil && *req.Make != "" {
		queryMods = append(queryMods, models.ReportVehicleSignalsEventsTrackingWhere.DeviceMake.EQ(*req.Make))
	}

	if req.Model != nil && *req.Model != "" {
		queryMods = append(queryMods, models.ReportVehicleSignalsEventsTrackingWhere.Model.EQ(*req.Model))
	}

	if req.DeviceDefinitionId != nil && *req.DeviceDefinitionId != "" {
		queryMods = append(queryMods, models.ReportVehicleSignalsEventsTrackingWhere.DeviceDefinitionID.EQ(*req.DeviceDefinitionId))
	}

	if req.Year != nil && *req.Year > 0 {
		queryMods = append(queryMods, models.ReportVehicleSignalsEventsTrackingWhere.Year.EQ(int(*req.Year)))
	}

	var eventProperties []*internalmodel.SignalsEvents
	err := models.ReportVehicleSignalsEventsTrackings(queryMods...).Bind(ctx, s.dbs().Reader, &eventProperties)

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error. "+err.Error())
	}

	queryAllMods := []qm.QueryMod{
		qm.Select(fmt.Sprintf("%s as name", keys[*req.Level]), "SUM(count) as total_count"),
		qm.GroupBy(keys[*req.Level]),
	}

	queryAllMods = append(queryAllMods, models.ReportVehicleSignalsEventsAllWhere.IntegrationID.EQ(req.IntegrationId))
	queryAllMods = append(queryAllMods, models.ReportVehicleSignalsEventsAllWhere.DateID.EQ(req.DateId))

	if req.PropertyId != nil && *req.PropertyId != "" {
		queryAllMods = append(queryAllMods, models.ReportVehicleSignalsEventsAllWhere.PropertyID.EQ(*req.PropertyId))
	}

	if req.Make != nil && *req.Make != "" {
		queryAllMods = append(queryAllMods, models.ReportVehicleSignalsEventsAllWhere.DeviceMake.EQ(*req.Make))
	}

	if req.Model != nil && *req.Model != "" {
		queryAllMods = append(queryAllMods, models.ReportVehicleSignalsEventsAllWhere.Model.EQ(*req.Model))
	}

	if req.DeviceDefinitionId != nil && *req.DeviceDefinitionId != "" {
		queryAllMods = append(queryAllMods, models.ReportVehicleSignalsEventsAllWhere.DeviceDefinitionID.EQ(*req.DeviceDefinitionId))
	}

	if req.Year != nil && *req.Year > 0 {
		queryAllMods = append(queryAllMods, models.ReportVehicleSignalsEventsAllWhere.Year.EQ(int(*req.Year)))
	}

	var allEvents []*internalmodel.SignalsEvents
	err = models.ReportVehicleSignalsEventsAlls(queryAllMods...).Bind(ctx, s.dbs().Reader, &allEvents)

	if err != nil {
		s.logger.Info().Msg("Running second query (failed)")
		return nil, status.Error(codes.Internal, "Internal error."+err.Error())
	}

	availableProperties, err := models.VehicleSignalsAvailableProperties().All(ctx, s.dbs().Reader)
	if err != nil {
		s.logger.Info().Msg("VehicleSignalsAvailableProperties (failed)")
		return nil, status.Error(codes.Internal, "Internal error."+err.Error())
	}

	summaryData, err := s.GetSummaryConnected(ctx, &pb.SummaryConnectedRequest{
		DateId:        req.DateId,
		IntegrationId: req.IntegrationId,
	})

	result := &pb.SignalResponse{}
	for _, event := range allEvents {
		requestCount := 0
		for _, eventProperty := range eventProperties {
			if strings.TrimSpace(eventProperty.Name) == strings.TrimSpace(event.Name) {
				requestCount = int(eventProperty.TotalCount)
				break
			}
		}

		if req.RemoveEmpty != nil && *req.RemoveEmpty {
			if requestCount == 0 && event.TotalCount == 0 {
				continue
			}
		}

		signalItemResponse := &pb.SignalItemResponse{
			Name:         strings.TrimSpace(event.Name),
			RequestCount: int32(requestCount),
			TotalCount:   int32(event.TotalCount),
		}

		for _, prop := range availableProperties {
			//s.logger.Info().Msgf("prop.ID => %s | event.Name => %s", prop.ID, signalItemResponse.Name)
			if strings.TrimSpace(prop.ID) == strings.TrimSpace(signalItemResponse.Name) {
				signalItemResponse.PowerTrainType = prop.PowerTrainType
				if len(signalItemResponse.PowerTrainType) > 0 {
					var powerTrainTypeTotalCount int32
					for _, powerTrainType := range signalItemResponse.PowerTrainType {

						if len(summaryData.PowerTrainTypeCountTimeframe) > 0 {
							for _, powerTrainTypeItem := range summaryData.PowerTrainTypeCountTimeframe {
								if strings.TrimSpace(powerTrainType) == strings.TrimSpace(powerTrainTypeItem.Type) {
									powerTrainTypeTotalCount += powerTrainTypeItem.Count
								}
							}
							signalItemResponse.TotalCount = powerTrainTypeTotalCount
						}
					}

				}
				break
			}
		}

		result.Items = append(result.Items, signalItemResponse)
	}

	return result, nil
}

func (s *userDeviceData) GetAvailableDates(ctx context.Context, _ *emptypb.Empty) (*pb.DateIdsResponse, error) {
	// raw query, project to list of strings

	query := `select date_id, integration_id from
(select date_id, integration_id
from device_data_api.report_vehicle_signals_events_tracking
group by date_id, integration_id) as dates
order by date_id desc`

	// need obj array
	var dateIDSlice []*internalmodel.DateIDItem

	err := queries.Raw(query).Bind(ctx, s.dbs().Reader, &dateIDSlice)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := &pb.DateIdsResponse{
		DateIds: make([]*pb.DateIdResponseItem, len(dateIDSlice)),
	}
	for i, item := range dateIDSlice {
		result.DateIds[i] = &pb.DateIdResponseItem{
			DateId:        item.DateID,
			IntegrationId: item.IntegrationID,
		}
	}

	return result, nil
}

func (s *userDeviceData) GetSummaryConnected(ctx context.Context, in *pb.SummaryConnectedRequest) (*pb.SummaryConnectedResponse, error) {
	allTimeCnt, err := models.ReportVehicleSignalsEventsSummaries(
		models.ReportVehicleSignalsEventsSummaryWhere.IntegrationID.EQ(in.IntegrationId),
	).Count(ctx, s.dbs().Reader)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := &pb.SummaryConnectedResponse{
		ConnectedAllTime:   allTimeCnt,
		ConnectedTimeframe: 0,
	}

	queryMods := []qm.QueryMod{
		qm.Select(models.ReportVehicleSignalsEventsSummaryColumns.IntegrationID, models.ReportVehicleSignalsEventsSummaryColumns.PowerTrainType, "SUM(count) as total_count", "SUM(device_definition_count) as total_device_definition_count"),
		qm.GroupBy(models.ReportVehicleSignalsEventsSummaryColumns.IntegrationID + "," + models.ReportVehicleSignalsEventsSummaryColumns.PowerTrainType),
	}

	queryMods = append(queryMods, models.ReportVehicleSignalsEventsSummaryWhere.IntegrationID.EQ(in.IntegrationId))
	queryMods = append(queryMods, models.ReportVehicleSignalsEventsSummaryWhere.DateID.EQ(in.DateId))

	var allEvents []*internalmodel.SignalsEventsUserDevices
	err = models.ReportVehicleSignalsEventsSummaries(queryMods...).Bind(ctx, s.dbs().Reader, &allEvents)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error."+err.Error())
	}

	if len(allEvents) == 0 {
		result.DateRange = "No Data found for Integration and Date"
		return result, nil
	}

	// build date object from in.DateId
	endDate, err := convertToDate(in.DateId)
	if err != nil {
		return nil, err
	}
	result.DateRange = endDate.Add(time.Hour*24*-7).Format(time.RFC1123) + " to " + endDate.Format(time.RFC1123)

	result.StartDate = endDate.Add(time.Hour * 24 * -7).Format(time.RFC1123)
	result.EndDate = endDate.Format(time.RFC1123)

	powerTrainTypeGroups := make(map[string][]*internalmodel.SignalsEventsUserDevices)
	connectedTimeframe := 0
	for _, item := range allEvents {
		connectedTimeframe += int(item.TotalCount)
		powerTrainTypeGroups[item.PowerTrainType] = append(powerTrainTypeGroups[item.PowerTrainType], item)
	}

	result.ConnectedTimeframe = int64(connectedTimeframe)

	for powerTrainType, group := range powerTrainTypeGroups {
		powerTrainTypeTimeframeCount := 0
		powerTrainTypeTimeframeDeviceDefinitionCount := 0
		for _, item := range group {
			powerTrainTypeTimeframeCount += int(item.TotalCount)
			powerTrainTypeTimeframeDeviceDefinitionCount += int(item.TotalDeviceDefinitionCount)
		}
		result.PowerTrainTypeCountTimeframe = append(result.PowerTrainTypeCountTimeframe, &pb.SummaryConnectedResponse_PowerTrainTypeConnectedResponse{
			Type:                  strings.TrimSpace(powerTrainType),
			Count:                 int32(powerTrainTypeTimeframeCount),
			DeviceDefinitionCount: int32(powerTrainTypeTimeframeDeviceDefinitionCount),
		})
	}

	return result, nil
}

func convertToDate(input string) (time.Time, error) {
	// Check if the input string is valid and has a length of 8 characters
	if len(input) != 8 {
		return time.Time{}, fmt.Errorf("invalid input: must be 8 characters long")
	}

	// Extract year, month, and day from the input string
	year := input[:4]
	month := input[4:6]
	day := input[6:8]

	// Parse the extracted parts into integers
	yearInt := 0
	monthInt := 0
	dayInt := 0

	_, err := fmt.Sscanf(year, "%d", &yearInt)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year: %v", err)
	}

	_, err = fmt.Sscanf(month, "%d", &monthInt)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid month: %v", err)
	}

	_, err = fmt.Sscanf(day, "%d", &dayInt)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day: %v", err)
	}

	// Create the date from the extracted parts
	date := time.Date(yearInt, time.Month(monthInt), dayInt, 0, 0, 0, 0, time.UTC)
	return date, nil
}

func isEmpty(currentString string) bool {
	return len(strings.TrimSpace(currentString)) == 0
}

package rpc

import (
	"context"
	"database/sql"
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

// todo need test for this

func (s *userDeviceData) GetUserDeviceData(ctx context.Context, req *pb.UserDeviceDataRequest) (*pb.UserDeviceDataResponse, error) {
	if req.UserDeviceId == "" || req.DeviceDefinitionId == "" {
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
	queryMods := []qm.QueryMod{
		qm.Select("property_id", "SUM(count) as total_count"),
		qm.GroupBy("property_id"),
	}

	queryMods = append(queryMods, models.ReportVehicleSignalsEventsTrackingWhere.IntegrationID.EQ(req.IntegrationId))
	queryMods = append(queryMods, models.ReportVehicleSignalsEventsTrackingWhere.DateID.EQ(req.DateId))

	if req.SignalName != nil && *req.SignalName != "" {
		queryMods = append(queryMods, models.ReportVehicleSignalsEventsTrackingWhere.PropertyID.EQ(*req.SignalName))
	}

	if req.MakeId != nil && *req.MakeId != "" {
		queryMods = append(queryMods, models.ReportVehicleSignalsEventsTrackingWhere.DeviceMakeID.EQ(*req.MakeId))
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
		qm.Select("property_id", "SUM(count) as total_count"),
		qm.GroupBy("property_id"),
	}

	queryAllMods = append(queryAllMods, models.ReportVehicleSignalsEventsAllWhere.IntegrationID.EQ(req.IntegrationId))
	queryAllMods = append(queryAllMods, models.ReportVehicleSignalsEventsAllWhere.DateID.EQ(req.DateId))

	if req.SignalName != nil && *req.SignalName != "" {
		queryAllMods = append(queryAllMods, models.ReportVehicleSignalsEventsAllWhere.PropertyID.EQ(*req.SignalName))
	}

	if req.MakeId != nil && *req.MakeId != "" {
		queryAllMods = append(queryAllMods, models.ReportVehicleSignalsEventsAllWhere.DeviceMakeID.EQ(*req.MakeId))
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
		return nil, status.Error(codes.Internal, "Internal error."+err.Error())
	}

	result := &pb.SignalResponse{}
	for _, event := range allEvents {
		requestCount := 0
		for _, eventProperty := range eventProperties {
			if eventProperty.PropertyID == event.PropertyID {
				requestCount = int(eventProperty.TotalCount)
				break
			}
		}
		result.Items = append(result.Items, &pb.SignalItemResponse{
			Property:     event.PropertyID,
			RequestCount: int32(requestCount),
			TotalCount:   int32(event.TotalCount),
		})
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

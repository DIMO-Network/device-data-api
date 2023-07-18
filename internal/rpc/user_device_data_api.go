package rpc

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/queries"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/pkg/errors"
	smartcar "github.com/smartcar/go-sdk"

	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/DIMO-Network/device-data-api/models"
	pb "github.com/DIMO-Network/device-data-api/pkg/grpc"
	"github.com/DIMO-Network/shared/db"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	queryEventProperty := qm.Where(
		models.ReportVehicleSignalsEventsTrackingColumns.IntegrationID+" = ?",
		req.IntegrationId,
		qm.Where(models.ReportVehicleSignalsEventsTrackingColumns.DateID+" = ?", req.DateId),
	)

	eventProperties, err := models.ReportVehicleSignalsEventsTrackings(queryEventProperty).All(ctx, s.dbs().Reader)

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	queryEvent := qm.Where(
		models.ReportVehicleSignalsEventsAllColumns.IntegrationID+" = ?",
		req.IntegrationId,
		qm.Where(models.ReportVehicleSignalsEventsAllColumns.DateID+" = ?", req.DateId),
	)

	events, err := models.ReportVehicleSignalsEventsAlls(queryEvent).All(ctx, s.dbs().Reader)

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	result := &pb.SignalResponse{}
	for _, event := range events {
		requestCount := 0
		for _, eventProperty := range eventProperties {
			if eventProperty.PropertyID == event.PropertyID {
				requestCount = eventProperty.Count
				break
			}
		}
		result.Items = append(result.Items, &pb.SignalItemResponse{
			Property:     event.PropertyID,
			RequestCount: int32(requestCount),
			TotalCount:   int32(event.Count),
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
	var dateIDSlice []*dateIDItem

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

type dateIDItem struct {
	DateID        string `boil:"date_id" json:"date_id" toml:"date_id" yaml:"date_id"`
	IntegrationID string `boil:"integration_id" json:"integration_id" toml:"integration_id" yaml:"integration_id"`
}

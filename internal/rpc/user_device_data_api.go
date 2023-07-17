package rpc

import (
	"context"
	"database/sql"

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
		Charging:             convertBoolPtr(ds.Charging),
		FuelPercentRemaining: convertFloatPtr(ds.FuelPercentRemaining),
		BatteryCapacity:      convertIntPtr(ds.BatteryCapacity),
		OilLevel:             convertFloatPtr(ds.OilLevel),
		Odometer:             convertFloatPtr(ds.Odometer),
		Latitude:             convertFloatPtr(ds.Latitude),
		Longitude:            convertFloatPtr(ds.Longitude),
		Range:                convertFloatPtr(ds.Range),
		StateOfCharge:        convertFloatPtr(ds.StateOfCharge),
		ChargeLimit:          convertFloatPtr(ds.ChargeLimit),
		RecordUpdatedAt:      convertToTimestamp(ds.RecordUpdatedAt),
		RecordCreatedAt:      convertToTimestamp(ds.RecordCreatedAt),
		TirePressure:         convertTirePressure(ds.TirePressure),
		BatteryVoltage:       convertFloatPtr(ds.BatteryVoltage),
		AmbientTemp:          convertFloatPtr(ds.AmbientTemp),
	}, nil
}

func convertBoolPtr(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
func convertFloatPtr(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
func convertIntPtr(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
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

	fromDate := req.FromDate.AsTime().Format("20060102")
	toDate := req.ToDate.AsTime().Format("20060102")

	queryEventProperty := qm.Where(
		models.ReportVehicleSignalsEventsTrackingColumns.IntegrationID+" = ?",
		req.IntegrationId,
		qm.WhereIn(models.ReportVehicleSignalsEventsTrackingColumns.DateID, []string{fromDate, toDate}),
	)

	eventProperties, err := models.ReportVehicleSignalsEventsTrackings(queryEventProperty).All(ctx, s.dbs().Reader)

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	queryEvent := qm.Where(
		models.ReportVehicleSignalsEventsAllColumns.IntegrationID+" = ?",
		req.IntegrationId,
		qm.WhereIn(models.ReportVehicleSignalsEventsAllColumns.DateID, []string{fromDate, toDate}),
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

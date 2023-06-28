package api

import (
	"context"
	"strconv"

	"sort"
	"time"

	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DIMO-Network/device-data-api/internal/constants"
	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/DIMO-Network/device-data-api/models"
	pb "github.com/DIMO-Network/device-data-api/pkg/grpc"
	"github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
	"github.com/volatiletech/null/v8"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func NewUserDeviceData(dbs func() *db.ReaderWriter, logger *zerolog.Logger, deviceDefSvc services.DeviceDefinitionsAPIService) pb.UseDeviceDataServiceServer {
	return &userDeviceData{dbs: dbs, logger: logger, deviceDefSvc: deviceDefSvc}
}

type userDeviceData struct {
	pb.UseDeviceDataServiceServer
	dbs          func() *db.ReaderWriter
	logger       *zerolog.Logger
	deviceDefSvc services.DeviceDefinitionsAPIService
}

func (s *userDeviceData) GetUserDeviceData(ctx context.Context, req *pb.UserDeviceDataRequest) (*pb.UserDeviceDataResponse, error) {
	deviceData, err := models.UserDeviceData(
		models.UserDeviceDatumWhere.UserDeviceID.EQ(req.UserDeviceId),
		models.UserDeviceDatumWhere.Signals.IsNotNull(),
		models.UserDeviceDatumWhere.UpdatedAt.GT(time.Now().Add(-14*24*time.Hour)),
	).All(ctx, s.dbs().Reader)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	if len(deviceData) == 0 {
		return nil, status.Error(codes.NotFound, "No status updates yet.")
	}

	ds := prepareDeviceStatusInformation(ctx, s.deviceDefSvc, deviceData, req.DeviceDefinitionId,
		null.StringFrom(req.DeviceStyleId), []int64{constants.NonLocationData, constants.CurrentLocation, constants.AllTimeLocation})

	return ds, nil
}

func (s *userDeviceData) GetSignals(ctx context.Context, req *pb.SignalRequest) (*pb.SignalResponse, error) {

	fromDate := req.FromDate.AsTime().Format("20060102")
	toDate := req.ToDate.AsTime().Format("20060102")

	query := qm.Where(
		models.ReportVehicleSignalsEventsPropertyColumns.IntegrationID+" = ?",
		req.IntegrationId,
		qm.WhereIn(models.ReportVehicleSignalsEventsPropertyColumns.DateID, []string{fromDate, toDate}),
	)

	events, err := models.ReportVehicleSignalsEventsProperties(query).All(ctx, s.dbs().Reader)

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	count, err := models.UserDeviceData(
		qm.Where("updated_at > ? AND updated_at < ?", req.FromDate, req.ToDate),
		qm.Where("integration_id = ?", req.IntegrationId),
	).Count(ctx, s.dbs().Reader)

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	result := &pb.SignalResponse{}

	for _, event := range events {
		result.Items = append(result.Items, &pb.SignalItemResponse{
			Property:     event.PropertyID,
			RequestCount: int32(event.Count),
			TotalCount:   int32(count),
		})
	}

	return result, nil
}

func prepareDeviceStatusInformation(ctx context.Context, ddSvc services.DeviceDefinitionsAPIService, deviceData models.UserDeviceDatumSlice, deviceDefinitionID string, deviceStyleID null.String, privilegeIDs []int64) *pb.UserDeviceDataResponse {
	ds := pb.UserDeviceDataResponse{}

	// set the record created date to most recent one
	for _, datum := range deviceData {
		if ds.RecordCreatedAt == nil || convertToUnixTimestamp(ds.RecordCreatedAt) < datum.CreatedAt.Unix() {
			ds.RecordCreatedAt = convertToTimestamp(datum.CreatedAt)
		}
	}

	var hasRage bool
	var hasFuelPercentRemaining bool
	// future: if time btw UpdateAt and timestamp > 7 days, ignore property

	// todo further refactor by passing in type for each option, then have switch in function below, can also refactor timestamp thing
	if slices.Contains(privilegeIDs, constants.NonLocationData) {
		charging := findMostRecentSignal(deviceData, "charging", false)
		if charging.Exists() {
			c := charging.Get("value").Bool()
			ds.Charging = c
		}
		fuelPercentRemaining := findMostRecentSignal(deviceData, "fuelPercentRemaining", false)
		if fuelPercentRemaining.Exists() {
			ts := fuelPercentRemaining.Get("timestamp").Time()
			if ds.RecordUpdatedAt == nil || convertToUnixTimestamp(ds.RecordUpdatedAt) < ts.Unix() {
				ds.RecordUpdatedAt = convertToTimestamp(ts)
			}
			f := fuelPercentRemaining.Get("value").Float()
			if f >= 0.01 {
				ds.FuelPercentRemaining = f
				hasFuelPercentRemaining = true
			}
		}
		batteryCapacity := findMostRecentSignal(deviceData, "batteryCapacity", false)
		if batteryCapacity.Exists() {
			b := batteryCapacity.Get("value").Int()
			ds.BatteryCapacity = b
		}
		oilLevel := findMostRecentSignal(deviceData, "oil", false)
		if oilLevel.Exists() {
			o := oilLevel.Get("value").Float()
			ds.OilLevel = o
		}
		stateOfCharge := findMostRecentSignal(deviceData, "soc", false)
		if stateOfCharge.Exists() {
			o := stateOfCharge.Get("value").Float()
			ds.StateOfCharge = o
		}
		chargeLimit := findMostRecentSignal(deviceData, "chargeLimit", false)
		if chargeLimit.Exists() {
			o := chargeLimit.Get("value").Float()
			ds.ChargeLimit = o
		}
		odometer := findMostRecentSignal(deviceData, "odometer", true)
		if odometer.Exists() {
			ts := odometer.Get("timestamp").Time()
			if ds.RecordUpdatedAt == nil || convertToUnixTimestamp(ds.RecordUpdatedAt) < ts.Unix() {
				ds.RecordUpdatedAt = convertToTimestamp(ts)
			}
			o := odometer.Get("value").Float()
			if shared.IsOdometerValid(o) {
				ds.Odometer = o
			}
		}
		rangeG := findMostRecentSignal(deviceData, "range", false)

		if rangeG.Exists() {
			r := rangeG.Get("value").Float()
			ds.Range = r
			hasRage = true
		}
		batteryVoltage := findMostRecentSignal(deviceData, "batteryVoltage", false)
		if batteryVoltage.Exists() {
			ts := batteryVoltage.Get("timestamp").Time()
			if ds.RecordUpdatedAt == nil || convertToUnixTimestamp(ds.RecordUpdatedAt) < ts.Unix() {
				ds.RecordUpdatedAt = convertToTimestamp(ts)
			}
			bv := batteryVoltage.Get("value").Float()
			ds.BatteryVoltage = bv
		}
		ambientTemp := findMostRecentSignal(deviceData, "ambientTemp", false)
		if ambientTemp.Exists() {
			at := ambientTemp.Get("value").Float()
			ds.AmbientTemp = at
		}
		// TirePressure
		tires := findMostRecentSignal(deviceData, "tires", false)
		if tires.Exists() {
			// weird thing here is in example payloads these are all ints, but the smartcar lib has as floats
			ds.TirePressure = &pb.TirePressureResponse{
				FrontLeft:  tires.Get("value").Get("frontLeft").Float(),
				FrontRight: tires.Get("value").Get("frontRight").Float(),
				BackLeft:   tires.Get("value").Get("backLeft").Float(),
				BackRight:  tires.Get("value").Get("backRight").Float(),
			}
		}
	}

	if slices.Contains(privilegeIDs, constants.CurrentLocation) || slices.Contains(privilegeIDs, constants.AllTimeLocation) {
		latitude := findMostRecentSignal(deviceData, "latitude", false)
		if latitude.Exists() {
			ts := latitude.Get("timestamp").Time()
			if ds.RecordUpdatedAt == nil || convertToUnixTimestamp(ds.RecordUpdatedAt) < ts.Unix() {
				ds.RecordUpdatedAt = convertToTimestamp(ts)
			}
			l := latitude.Get("value").Float()
			ds.Latitude = l
		}
		longitude := findMostRecentSignal(deviceData, "longitude", false)
		if longitude.Exists() {
			l := longitude.Get("value").Float()
			ds.Longitude = l
		}
	}

	if !hasRage && hasFuelPercentRemaining {
		calcRange, err := calculateRange(ctx, ddSvc, deviceDefinitionID, deviceStyleID, ds.FuelPercentRemaining)
		if err == nil {
			ds.Range = *calcRange
		}
	}

	return &ds
}

func convertToTimestamp(goTime time.Time) *timestamppb.Timestamp {
	timestamp := timestamppb.New(goTime)
	return timestamp
}

func convertToUnixTimestamp(timestamp *timestamppb.Timestamp) int64 {
	goTime := timestamp.AsTime()
	unixTimestamp := goTime.Unix()
	return unixTimestamp
}

// findMostRecentSignal finds the highest value float instead of most recent, eg. for odometer
func findMostRecentSignal(udd models.UserDeviceDatumSlice, path string, highestFloat bool) gjson.Result {
	// todo write test
	if len(udd) == 0 {
		return gjson.Result{}
	}
	if len(udd) > 1 {
		if highestFloat {
			sortBySignalValueDesc(udd, path)
		} else {
			sortBySignalTimestampDesc(udd, path)
		}
	}
	return gjson.GetBytes(udd[0].Signals.JSON, path)
}

// calculateRange returns the current estimated range based on fuel tank capacity, mpg, and fuelPercentRemaining and returns it in Kilometers
func calculateRange(ctx context.Context, ddSvc services.DeviceDefinitionsAPIService, deviceDefinitionID string, deviceStyleID null.String, fuelPercentRemaining float64) (*float64, error) {
	if fuelPercentRemaining <= 0.01 {
		return nil, status.Error(codes.Internal, "fuelPercentRemaining lt 0.01 so cannot calculate range")
	}

	dd, err := ddSvc.GetDeviceDefinitionsByIDs(ctx, []string{deviceDefinitionID})

	if err != nil {
		return nil, err
	}

	rangeData := getActualDeviceDefinitionMetadataValues(dd[0], deviceStyleID)

	// calculate, convert to Km
	if rangeData.FuelTankCapGal > 0 && rangeData.Mpg > 0 {
		fuelTankAtGal := rangeData.FuelTankCapGal * fuelPercentRemaining
		rangeMiles := rangeData.Mpg * fuelTankAtGal
		rangeKm := 1.60934 * rangeMiles
		return &rangeKm, nil
	}

	return nil, nil
}

// sortBySignalValueDesc Sort user device data so the highest value is first
func sortBySignalValueDesc(udd models.UserDeviceDatumSlice, path string) {
	sort.Slice(udd, func(i, j int) bool {
		fpri := gjson.GetBytes(udd[i].Signals.JSON, path+".value")
		fprj := gjson.GetBytes(udd[j].Signals.JSON, path+".value")
		// if one has it and the other does not, makes no difference
		if fpri.Exists() && !fprj.Exists() {
			return true
		} else if !fpri.Exists() && fprj.Exists() {
			return false
		}
		return fprj.Float() < fpri.Float()
	})
}

// sortBySignalTimestampDesc Sort user device data so the most recent timestamp is first
func sortBySignalTimestampDesc(udd models.UserDeviceDatumSlice, path string) {
	sort.Slice(udd, func(i, j int) bool {
		fpri := gjson.GetBytes(udd[i].Signals.JSON, path+".timestamp")
		fprj := gjson.GetBytes(udd[j].Signals.JSON, path+".timestamp")
		// if one has it and the other does not, makes no difference
		if fpri.Exists() && !fprj.Exists() {
			return true
		} else if !fpri.Exists() && fprj.Exists() {
			return false
		}
		return fprj.Time().Unix() < fpri.Time().Unix()
	})
}

func getActualDeviceDefinitionMetadataValues(dd *grpc.GetDeviceDefinitionItemResponse, deviceStyleID null.String) *DeviceDefinitionRange {

	var fuelTankCapGal, mpg, mpgHwy float64 = 0, 0, 0

	var metadata []*grpc.DeviceTypeAttribute

	if !deviceStyleID.IsZero() {
		for _, style := range dd.DeviceStyles {
			if style.Id == deviceStyleID.String {
				metadata = style.DeviceAttributes
				break
			}
		}
	}

	if len(metadata) == 0 && dd != nil && dd.DeviceAttributes != nil {
		metadata = dd.DeviceAttributes
	}

	for _, attr := range metadata {
		switch DeviceAttributeType(attr.Name) {
		case FuelTankCapacityGal:
			if v, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				fuelTankCapGal = v
			}
		case Mpg:
			if v, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				mpg = v
			}
		case MpgHighway:
			if v, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				mpgHwy = v
			}
		}
	}

	return &DeviceDefinitionRange{
		FuelTankCapGal: fuelTankCapGal,
		Mpg:            mpg,
		MpgHwy:         mpgHwy,
	}
}

type DeviceAttributeType string

const (
	Mpg                 DeviceAttributeType = "mpg"
	FuelTankCapacityGal DeviceAttributeType = "fuel_tank_capacity_gal"
	MpgHighway          DeviceAttributeType = "mpg_highway"
)

type DeviceDefinitionRange struct {
	FuelTankCapGal float64 `json:"fuel_tank_capacity_gal"`
	Mpg            float64 `json:"mpg"`
	MpgHwy         float64 `json:"mpg_highway"`
}

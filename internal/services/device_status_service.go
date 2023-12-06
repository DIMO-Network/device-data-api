package services

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/response"
	"github.com/DIMO-Network/device-data-api/models"
	"github.com/DIMO-Network/shared"
	smartcar "github.com/smartcar/go-sdk"
	"github.com/tidwall/gjson"
	"github.com/volatiletech/null/v8"
	"golang.org/x/exp/slices"
)

//go:generate mockgen -source device_status_service.go -destination mocks/device_status_service_mock.go
type deviceStatusService struct {
	ddSvc DeviceDefinitionsAPIService
}

type DeviceStatusService interface {
	PrepareDeviceStatusInformation(ctx context.Context, deviceData models.UserDeviceDatumSlice, deviceDefinitionID string, deviceStyleID *string, privilegeIDs []int64) response.DeviceSnapshot
	PrepareDeviceStatusInformationV2(ctx context.Context, deviceData models.UserDeviceDatumSlice, deviceDefinitionID string, deviceStyleID *string, privilegeIDs []int64) response.Device
	CalculateRange(ctx context.Context, deviceDefinitionID string, deviceStyleID *string, fuelPercentRemaining float64) (*float64, error)
}

func NewDeviceStatusService(deviceDefinitionsSvc DeviceDefinitionsAPIService) DeviceStatusService {
	return &deviceStatusService{
		ddSvc: deviceDefinitionsSvc,
	}
}

const (
	NonLocationData int64 = 1
	Commands        int64 = 2
	CurrentLocation int64 = 3
	AllTimeLocation int64 = 4
)

func (dss *deviceStatusService) PrepareDeviceStatusInformation(ctx context.Context, deviceData models.UserDeviceDatumSlice, deviceDefinitionID string, deviceStyleID *string, privilegeIDs []int64) response.DeviceSnapshot {
	ds := response.DeviceSnapshot{}

	// set the record created date to most recent one
	for _, datum := range deviceData {
		if ds.RecordCreatedAt == nil || ds.RecordCreatedAt.Unix() < datum.CreatedAt.Unix() {
			ds.RecordCreatedAt = &datum.CreatedAt
		}
	}
	// future: if time btw UpdateAt and timestamp > 7 days, ignore property

	// todo further refactor by passing in type for each option, then have switch in function below, can also refactor timestamp thing
	if slices.Contains(privilegeIDs, NonLocationData) {
		charging := findMostRecentSignal(deviceData, "charging", false)
		if charging.Exists() {
			c := charging.Get("value").Bool()
			ds.Charging = &c
		}
		fuelPercentRemaining := findMostRecentSignal(deviceData, "fuelPercentRemaining", false)
		if fuelPercentRemaining.Exists() {
			ts := fuelPercentRemaining.Get("timestamp").Time()
			if ds.RecordUpdatedAt == nil || ds.RecordUpdatedAt.Unix() < ts.Unix() {
				ds.RecordUpdatedAt = &ts
			}
			f := fuelPercentRemaining.Get("value").Float()
			if f >= 0.01 {
				ds.FuelPercentRemaining = &f
			}
		}
		batteryCapacity := findMostRecentSignal(deviceData, "batteryCapacity", false)
		if batteryCapacity.Exists() {
			b := batteryCapacity.Get("value").Int()
			ds.BatteryCapacity = &b
		}
		oilLevel := findMostRecentSignal(deviceData, "oil", false)
		if oilLevel.Exists() {
			o := oilLevel.Get("value").Float()
			ds.OilLevel = &o
		}
		stateOfCharge := findMostRecentSignal(deviceData, "soc", false)
		if stateOfCharge.Exists() {
			o := stateOfCharge.Get("value").Float()
			ds.StateOfCharge = &o
		}
		chargeLimit := findMostRecentSignal(deviceData, "chargeLimit", false)
		if chargeLimit.Exists() {
			o := chargeLimit.Get("value").Float()
			ds.ChargeLimit = &o
		}
		odometer := findMostRecentSignal(deviceData, "odometer", true)
		if odometer.Exists() {
			ts := odometer.Get("timestamp").Time()
			if ds.RecordUpdatedAt == nil || ds.RecordUpdatedAt.Unix() < ts.Unix() {
				ds.RecordUpdatedAt = &ts
			}
			o := odometer.Get("value").Float()
			if shared.IsOdometerValid(o) {
				ds.Odometer = &o
			}
		}
		rangeG := findMostRecentSignal(deviceData, "range", false)
		if rangeG.Exists() {
			r := rangeG.Get("value").Float()
			ds.Range = &r
		}
		batteryVoltage := findMostRecentSignal(deviceData, "batteryVoltage", false)
		if batteryVoltage.Exists() {
			ts := batteryVoltage.Get("timestamp").Time()
			if ds.RecordUpdatedAt == nil || ds.RecordUpdatedAt.Unix() < ts.Unix() {
				ds.RecordUpdatedAt = &ts
			}
			bv := batteryVoltage.Get("value").Float()
			ds.BatteryVoltage = &bv
		}
		ambientTemp := findMostRecentSignal(deviceData, "ambientTemp", false)
		if ambientTemp.Exists() {
			at := ambientTemp.Get("value").Float()
			ds.AmbientTemp = &at
		}
		// TirePressure
		tires := findMostRecentSignal(deviceData, "tires", false)
		if tires.Exists() {
			// weird thing here is in example payloads these are all ints, but the smartcar lib has as floats
			ds.TirePressure = &smartcar.TirePressure{
				FrontLeft:  tires.Get("value").Get("frontLeft").Float(),
				FrontRight: tires.Get("value").Get("frontRight").Float(),
				BackLeft:   tires.Get("value").Get("backLeft").Float(),
				BackRight:  tires.Get("value").Get("backRight").Float(),
			}
		}
	}

	if slices.Contains(privilegeIDs, CurrentLocation) || slices.Contains(privilegeIDs, AllTimeLocation) {
		latitude := findMostRecentSignal(deviceData, "latitude", false)
		if latitude.Exists() {
			ts := latitude.Get("timestamp").Time()
			if ds.RecordUpdatedAt == nil || ds.RecordUpdatedAt.Unix() < ts.Unix() {
				ds.RecordUpdatedAt = &ts
			}
			l := latitude.Get("value").Float()
			ds.Latitude = &l
		}
		longitude := findMostRecentSignal(deviceData, "longitude", false)
		if longitude.Exists() {
			l := longitude.Get("value").Float()
			ds.Longitude = &l
		}
	}

	if ds.Range == nil && ds.FuelPercentRemaining != nil {
		calcRange, err := dss.CalculateRange(ctx, deviceDefinitionID, deviceStyleID, *ds.FuelPercentRemaining)
		if err == nil {
			ds.Range = calcRange
		}
	}

	return ds
}

func (dss *deviceStatusService) PrepareDeviceStatusInformationV2(ctx context.Context, deviceData models.UserDeviceDatumSlice, deviceDefinitionID string, deviceStyleID *string, privilegeIDs []int64) response.Device {
	deviceSnapshot := response.Device{}

	// set the record created date to most recent one
	for _, datum := range deviceData {
		if deviceSnapshot.RecordCreatedAt == nil || deviceSnapshot.RecordCreatedAt.Unix() < datum.CreatedAt.Unix() {
			deviceSnapshot.RecordCreatedAt = &datum.CreatedAt
		}
	}
	// future: if time btw UpdateAt and timestamp > 7 days, ignore property

	// todo further refactor by passing in type for each option, then have switch in function below, can also refactor timestamp thing
	var status response.Status
	if slices.Contains(privilegeIDs, NonLocationData) {
		charging := findMostRecentSignal(deviceData, "charging", false)
		if charging.Exists() {
			c := charging.Get("value").Bool()
			status.PowerTrain.TractionBattery.Charging.IsCharging = null.BoolFrom(c)
		}
		fuelPercentRemaining := findMostRecentSignal(deviceData, "fuelPercentRemaining", false)
		if fuelPercentRemaining.Exists() {
			ts := fuelPercentRemaining.Get("timestamp").Time()
			if deviceSnapshot.RecordUpdatedAt == nil || deviceSnapshot.RecordUpdatedAt.Unix() < ts.Unix() {
				deviceSnapshot.RecordUpdatedAt = &ts
			}
			f := fuelPercentRemaining.Get("value").Float()
			if f >= 0.01 {
				status.PowerTrain.FuelSystem.Level = null.Float64From(f)
			}
		}
		batteryCapacity := findMostRecentSignal(deviceData, "batteryCapacity", false)
		if batteryCapacity.Exists() {
			b := batteryCapacity.Get("value").Float()
			status.PowerTrain.TractionBattery.GrossCapacity = null.Float64From(b)
			status.PowerTrain.TractionBattery.GrossCapacity = null.Float64From(b)
		}
		oilLevel := findMostRecentSignal(deviceData, "oil", false)
		if oilLevel.Exists() {
			o := oilLevel.Get("value").Float()
			if o > 0.75 {
				status.PowerTrain.CombustionEngine.EngineOilLevel = null.StringFrom("CRITICALLY_HIGH")
			} else if o >= 0.5 {
				status.PowerTrain.CombustionEngine.EngineOilLevel = null.StringFrom("NORMAL")
			} else if o >= 0.25 {
				status.PowerTrain.CombustionEngine.EngineOilLevel = null.StringFrom("LOW_NORMAL")
			} else if o > 0 {
				status.PowerTrain.CombustionEngine.EngineOilLevel = null.StringFrom("CRITICALLY_LOW")
			}

		}
		stateOfCharge := findMostRecentSignal(deviceData, "soc", false)
		if stateOfCharge.Exists() {
			o := stateOfCharge.Get("value").Float()
			status.PowerTrain.TractionBattery.StateOfCharge.Displayed = null.Float64From(o)
			status.PowerTrain.TractionBattery.StateOfCharge.Current = null.Float64From(o)
		}
		chargeLimit := findMostRecentSignal(deviceData, "chargeLimit", false)
		if chargeLimit.Exists() {
			o := chargeLimit.Get("value").Float()
			status.PowerTrain.TractionBattery.Charging.ChargeLimit = null.Float64From(o)
		}
		odometer := findMostRecentSignal(deviceData, "odometer", true)
		if odometer.Exists() {
			ts := odometer.Get("timestamp").Time()
			if deviceSnapshot.RecordUpdatedAt == nil || deviceSnapshot.RecordUpdatedAt.Unix() < ts.Unix() {
				deviceSnapshot.RecordUpdatedAt = &ts
			}
			o := odometer.Get("value").Float()
			if shared.IsOdometerValid(o) {
				status.TravelledDistance = null.Float64From(o)
				status.PowerTrain.Transmission.TravelledDistance = null.Float64From(o)
			}
		}
		rangeG := findMostRecentSignal(deviceData, "range", false)
		if rangeG.Exists() {
			r := rangeG.Get("value").Float()
			status.PowerTrain.Range = null.Float64From(r)
			status.PowerTrain.FuelSystem.Range = null.Float64From(r)
			status.PowerTrain.TractionBattery.Range = null.Float64From(r)
		}
		batteryVoltage := findMostRecentSignal(deviceData, "batteryVoltage", false)
		if batteryVoltage.Exists() {
			ts := batteryVoltage.Get("timestamp").Time()
			if deviceSnapshot.RecordUpdatedAt == nil || deviceSnapshot.RecordUpdatedAt.Unix() < ts.Unix() {
				deviceSnapshot.RecordUpdatedAt = &ts
			}
			bv := batteryVoltage.Get("value").Float()
			status.LowVoltageBattery.CurrentVoltage = null.Float64From(bv)
		}
		ambientTemp := findMostRecentSignal(deviceData, "ambientTemp", false)
		if ambientTemp.Exists() {
			at := ambientTemp.Get("value").Float()
			status.Exterior.AirTemperature = null.Float64From(at)
		}
		// TirePressure
		tires := findMostRecentSignal(deviceData, "tires", false)
		if tires.Exists() {
			// weird thing here is in example payloads these are all ints, but the smartcar lib has as floats
			fl := tires.Get("value").Get("frontLeft").Float()
			fr := tires.Get("value").Get("frontRight").Float()
			bl := tires.Get("value").Get("backLeft").Float()
			br := tires.Get("value").Get("backRight").Float()
			status.Chasis.Axle.Row1.Wheel.Left.Tire.Pressure = null.Float64From(fl)
			status.Chasis.Axle.Row1.Wheel.Right.Tire.Pressure = null.Float64From(fr)
			status.Chasis.Axle.Row2.Wheel.Left.Tire.Pressure = null.Float64From(bl)
			status.Chasis.Axle.Row2.Wheel.Right.Tire.Pressure = null.Float64From(br)

		}
	}

	if slices.Contains(privilegeIDs, CurrentLocation) || slices.Contains(privilegeIDs, AllTimeLocation) {
		latitude := findMostRecentSignal(deviceData, "latitude", false)
		if latitude.Exists() {
			ts := latitude.Get("timestamp").Time()
			if deviceSnapshot.RecordUpdatedAt == nil || deviceSnapshot.RecordUpdatedAt.Unix() < ts.Unix() {
				deviceSnapshot.RecordUpdatedAt = &ts
			}
			l := latitude.Get("value").Float()
			status.CurrentLocation.Latitude = null.Float64From(l)
			status.CurrentLocation.Timestamp = null.StringFrom(ts.Format(time.RFC3339))
		}
		longitude := findMostRecentSignal(deviceData, "longitude", false)
		if longitude.Exists() {
			l := longitude.Get("value").Float()
			status.CurrentLocation.Longitude = null.Float64From(l)
		}
	}

	if status.PowerTrain.TractionBattery.Range.IsZero() && !status.PowerTrain.FuelSystem.Level.IsZero() {
		calcRange, err := dss.CalculateRange(ctx, deviceDefinitionID, deviceStyleID, status.PowerTrain.FuelSystem.Level.Float64)
		if err == nil {
			status.PowerTrain.TractionBattery.Range = null.Float64From(*calcRange)
		}
	}

	deviceSnapshot.Status = append(deviceSnapshot.Status, status)

	return deviceSnapshot
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

// CalculateRange returns the current estimated range based on fuel tank capacity, mpg, and fuelPercentRemaining and returns it in Kilometers
func (dss *deviceStatusService) CalculateRange(ctx context.Context, deviceDefinitionID string, deviceStyleID *string, fuelPercentRemaining float64) (*float64, error) {
	if fuelPercentRemaining <= 0.01 {
		return nil, fmt.Errorf("fuelPercentRemaining lt 0.01 so cannot calculate range")
	}

	dd, err := dss.ddSvc.GetDeviceDefinitionByID(ctx, deviceDefinitionID)

	if err != nil {
		return nil, shared.GrpcErrorToFiber(err, "deviceDefSvc error getting definition id: "+deviceDefinitionID)
	}

	rangeData := GetActualDeviceDefinitionMetadataValues(dd, deviceStyleID)

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

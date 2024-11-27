package services

import (
	"context"
	"fmt"
	"sort"

	"github.com/DIMO-Network/device-data-api/internal/response"
	"github.com/DIMO-Network/device-data-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/privileges"
	smartcar "github.com/smartcar/go-sdk"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/slices"
)

//go:generate mockgen -source device_status_service.go -destination mocks/device_status_service_mock.go
type deviceStatusService struct {
	ddSvc DeviceDefinitionsAPIService
}

type DeviceStatusService interface {
	PrepareDeviceStatusInformation(ctx context.Context, deviceData models.UserDeviceDatumSlice, deviceDefinitionID string, deviceStyleID *string, privilegeIDs []privileges.Privilege) response.DeviceSnapshot
	CalculateRange(ctx context.Context, deviceDefinitionID string, deviceStyleID *string, fuelPercentRemaining float64) (*float64, error)
}

func NewDeviceStatusService(deviceDefinitionsSvc DeviceDefinitionsAPIService) DeviceStatusService {
	return &deviceStatusService{
		ddSvc: deviceDefinitionsSvc,
	}
}

func (dss *deviceStatusService) PrepareDeviceStatusInformation(ctx context.Context, deviceData models.UserDeviceDatumSlice, deviceDefinitionID string, deviceStyleID *string, privilegeIDs []privileges.Privilege) response.DeviceSnapshot {
	ds := response.DeviceSnapshot{}

	// set the record created date to most recent one
	for _, datum := range deviceData {
		if ds.RecordCreatedAt == nil || ds.RecordCreatedAt.Unix() < datum.CreatedAt.Unix() {
			ds.RecordCreatedAt = &datum.CreatedAt
		}
	}
	// future: if time btw UpdateAt and timestamp > 7 days, ignore property

	// todo further refactor by passing in type for each option, then have switch in function below, can also refactor timestamp thing
	if slices.Contains(privilegeIDs, privileges.VehicleNonLocationData) {
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

	if slices.Contains(privilegeIDs, privileges.VehicleCurrentLocation) || slices.Contains(privilegeIDs, privileges.VehicleAllTimeLocation) {
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

	dd, err := dss.ddSvc.GetDeviceDefinitionBySlug(ctx, deviceDefinitionID)
	if err != nil {
		return nil, shared.GrpcErrorToFiber(err, "deviceDefSvc error getting definition id: "+deviceDefinitionID)
	}
	// want the decimal form of the percentage for this calculation
	if fuelPercentRemaining > 1 {
		fuelPercentRemaining = fuelPercentRemaining / 100
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

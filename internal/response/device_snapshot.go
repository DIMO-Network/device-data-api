package response

import (
	"time"

	smartcar "github.com/smartcar/go-sdk"
	"github.com/volatiletech/null/v8"
)

// DeviceSnapshot is the response object for device status endpoint
// https://docs.google.com/document/d/1DYzzTOR9WA6WJNoBnwpKOoxfmrVwPWNLv0x0MkjIAqY/edit#heading=h.dnp7xngl47bw
type DeviceSnapshot struct {
	Charging             *bool                  `json:"charging,omitempty"`
	FuelPercentRemaining *float64               `json:"fuelPercentRemaining,omitempty"`
	BatteryCapacity      *int64                 `json:"batteryCapacity,omitempty"`
	OilLevel             *float64               `json:"oil,omitempty"`
	Odometer             *float64               `json:"odometer,omitempty"`
	Latitude             *float64               `json:"latitude,omitempty"`
	Longitude            *float64               `json:"longitude,omitempty"`
	Range                *float64               `json:"range,omitempty"`
	StateOfCharge        *float64               `json:"soc,omitempty"`
	ChargeLimit          *float64               `json:"chargeLimit,omitempty"`
	RecordUpdatedAt      *time.Time             `json:"recordUpdatedAt,omitempty"`
	RecordCreatedAt      *time.Time             `json:"recordCreatedAt,omitempty"`
	TirePressure         *smartcar.TirePressure `json:"tirePressure,omitempty"`
	BatteryVoltage       *float64               `json:"batteryVoltage,omitempty"`
	AmbientTemp          *float64               `json:"ambientTemp,omitempty"`
}

type Device struct {
	// VSS status data
	Status []Status `json:"status,omitempty"`
	// Other
	Misc            map[string]interface{} `json:"misc,omitempty"`
	RecordUpdatedAt *time.Time             `json:"recordUpdatedAt,omitempty"`
	RecordCreatedAt *time.Time             `json:"recordCreatedAt,omitempty"`
}

type Status struct {
	// Odometer reading, total distance travelled during the lifetime of the vehicle.
	TravelledDistance null.Float64 `json:"travelledDistance,omitempty"`
	// Vehicle speed.
	Speed           null.Float64 `json:"speed,omitempty"`
	CurrentLocation struct {
		// Current altitude relative to WGS 84 reference ellipsoid, as measured at the position of GNSS receiver antenna.
		Altitude null.Float64 `json:"altitude,omitempty"`
		// Current latitude of vehicle in WGS 84 geodetic coordinates, as measured at the position of GNSS receiver antenna.
		Latitude null.Float64 `json:"latitude,omitempty"`
		// Current longitude of vehicle in WGS 84 geodetic coordinates, as measured at the position of GNSS receiver antenna.
		Longitude null.Float64 `json:"longitude,omitempty"`
		// Timestamp from GNSS system for current location, formatted according to ISO 8601 with UTC time zone.
		Timestamp null.String `json:"timestamp,omitempty"`
	} `json:"currentLocation"`
	// Information about exterior measured by vehicle.
	Exterior struct {
		// Air temperature outside the vehicle.
		AirTemperature null.Float64 `json:"airTemperature,omitempty"`
	} `json:"exterior,omitempty"`

	// Signals related to low voltage battery.
	LowVoltageBattery struct {
		// Current Voltage of the low voltage battery.
		CurrentVoltage null.Float64 `json:"currentVoltage,omitempty"`
	} `json:"lowVoltageBattery,omitempty"`

	// OBD data.
	OBD struct {
		// PID 33 - Barometric pressure
		BarometricPressure null.Float64 `json:"barometricPressure,omitempty"`
		// PID 04 - Engine load in percent - 0 = no load, 100 = full load
		EngineLoad null.Float64 `json:"engineLoad,omitempty"`
		// PID 0C - Engine speed measured as rotations per minute
		EngineSpeed null.Float64 `json:"engineSpeed,omitempty"`
		// PID 51 - Fuel type
		FuelType null.String `json:"fuelType,omitempty"`
		// PID 0F - Intake temperature
		IntakeTemp null.Float64 `json:"intakeTemp,omitempty"`
		// PID 1F - Engine run time
		RunTime null.Float64 `json:"runTime,omitempty"`
		// PID 11 - Throttle position - 0 = closed throttle, 100 = open throttle
		ThrottlePosition null.Float64 `json:"throttlePosition,omitempty"`
	} `json:"obd,omitempty"`

	// Powertrain data for battery management, etc.
	PowerTrain struct {
		// Battery Management data.
		TractionBattery struct {
			// Gross capacity of the battery.
			GrossCapacity null.Float64 `json:"grossCapacity,omitempty"`
			// Total net capacity of the battery considering aging.
			NetCapacity null.Float64 `json:"netCapacity,omitempty"`
			// Properties related to battery charging.
			Charging struct {
				// Target charge limit (state of charge) for battery.
				ChargeLimit null.Float64 `json:"chargeLimit,omitempty"`
				// True if charging is ongoing. Charging is considered to be ongoing if energy is flowing from charger to vehicle.
				IsCharging null.Bool `json:"isCharging,omitempty"`
			} `json:"charging,omitempty"`
			// Remaining range in meters using only battery.
			Range null.Float64 `json:"range,omitempty"`
			// Information on the state of charge of the vehicle's high voltage battery.
			StateOfCharge struct {
				// Physical state of charge of the high voltage battery, relative to net capacity. This is not necessarily the state of charge being displayed to the customer.
				Current null.Float64 `json:"current,omitempty"`
				// State of charge displayed to the customer.
				Displayed null.Float64 `json:"displayed,omitempty"`
			} `json:"stateOfCharge,omitempty"`
		} `json:"tractionBattery,omitempty"`

		// Engine-specific data, stopping at the bell housing.
		CombustionEngine struct {
			// Engine coolant temperature.
			ECT null.Float64 `json:"ect,omitempty"`
			// Engine oil level. Must be one of: ['CRITICALLY_LOW', 'LOW', 'NORMAL', 'HIGH', 'CRITICALLY_HIGH']
			EngineOilLevel null.String `json:"engineOilLevel,omitempty"`
			// Engine speed measured as rotations per minute.
			Speed null.Float64 `json:"speed,omitempty"`
			// Current throttle position.
			TPS null.Float64 `json:"tps,omitempty"`
		} `json:"combustionEngine,omitempty"`

		// Transmission-specific data, stopping at the drive shafts.
		Transmission struct {
			// Odometer reading, total distance travelled during the lifetime of the transmission.
			TravelledDistance null.Float64 `json:"travelledDistance,omitempty"`
		} `json:"transmission,omitempty"`

		// Defines the powertrain type of the vehicle. For vehicles with a combustion engine (including hybrids) more detailed information on fuels supported can be found in FuelSystem.SupportedFuelTypes and FuelSystem.SupportedFuels. Must be one of: ['COMBUSTION', 'HYBRID', 'ELECTRIC']
		Type null.String `json:"type,omitempty"`

		// Fuel system data.
		FuelSystem struct {
			// Detailed information on fuels supported by the vehicle. Identifiers originating from DIN EN 16942:2021-08, appendix B, with additional suffix for octane (RON) where relevant. RON 95 is sometimes referred to as Super, RON 98 as Super Plus. Must be one of: ['E5_95', 'E5_98', 'E10_95', 'E10_98', 'E85', 'B7', 'B10', 'B20', 'B30', 'B100', 'XTL', 'LPG', 'CNG', 'LNG', 'H2', 'OTHER']
			SupportedFuel []null.String `json:"supportedFuel,omitempty"`
			// High level information of fuel types supported If a vehicle also has an electric drivetrain (e.g. hybrid) that will be obvious from the PowerTrain.Type signal. Must be one of: ['GASOLINE', 'DIESEL', 'E85', 'LPG', 'CNG', 'LNG', 'H2', 'OTHER']
			SupportedFuelTypes []null.String `json:"supportedFuelTypes,omitempty"`
			// Level in fuel tank as percent of capacity. 0 = empty. 100 = full.
			Level null.Float64 `json:"level,omitempty"`
			// Remaining range in meters using only liquid fuel.
			Range null.Float64 `json:"range,omitempty"`
		} `json:"fuelSystem,omitempty"`

		// Remaining range in meters using all energy sources available in the vehicle.
		Range null.Float64 `json:"range,omitempty"`
	} `json:"fuelSystem,omitempty"`

	// Attributes that identify a vehicle.
	VehicleIdentification struct {
		// Vehicle brand or manufacturer.
		Brand null.String `json:"brand,omitempty"`
		// Vehicle model.
		Model null.String `json:"model,omitempty"`
		// 17-character Vehicle Identification Number (VIN) as defined by ISO 3779.
		VIN null.String `json:"vin,omitempty"`
		// Model year of the vehicle.
		Year null.Int64 `json:"year,omitempty"`
	} `json:"vehicleIdentification,omitempty"`

	// All data concerning steering, suspension, wheels, and brakes.
	Chasis struct {
		// Axle signals
		Axle struct {
			Row1 struct {
				// Wheel signals for axle
				Wheel struct {
					// Wheel signals for axle
					Left struct {
						// Tire signals for wheel.
						Tire struct {
							// Tire pressure in kilo-Pascal.
							Pressure null.Float64 `json:"pressure,omitempty"`
						} `json:"left,omitempty"`
					} `json:"left,omitempty"`
					// Wheel signals for axle
					Right struct {
						// Tire signals for wheel.
						Tire struct {
							// Tire pressure in kilo-Pascal.
							Pressure null.Float64 `json:"pressure,omitempty"`
						} `json:"right,omitempty"`
					} `json:"right,omitempty"`
				} `json:"wheel,omitempty"`
			} `json:"row1,omitempty"`
			Row2 struct {
				// Wheel signals for axle
				Wheel struct {
					// Wheel signals for axle
					Left struct {
						// Tire signals for wheel.
						Tire struct {
							// Tire pressure in kilo-Pascal.
							Pressure null.Float64 `json:"pressure,omitempty"`
						} `json:"left,omitempty"`
					} `json:"left,omitempty"`
					// Wheel signals for axle
					Right struct {
						// Tire signals for wheel.
						Tire struct {
							// Tire pressure in kilo-Pascal.
							Pressure null.Float64 `json:"pressure,omitempty"`
						} `json:"right,omitempty"`
					} `json:"right,omitempty"`
				} `json:"wheel,omitempty"`
			} `json:"row2,omitempty"`
		} `json:"axle,omitempty"`
	} `json:"chassis,omitempty"`
}

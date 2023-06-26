-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

SET search_path = device_data_api, public;

CREATE TABLE IF NOT EXISTS vehicle_signals_tracking_properties
(
    id char(27) PRIMARY KEY,
    name varchar(200) not null,
    created_at     timestamptz not null default current_timestamp,
    updated_at     timestamptz not null default current_timestamp
    );

INSERT INTO vehicle_signals_tracking_properties(Id, name)
VALUES ('maf', 'maf'),
       ('vin', 'vin'),
       ('speed', 'speed'),
       ('runTime', 'runTime'),
       ('altitude', 'altitude'),
       ('latitude', 'latitude'),
       ('longitude', 'longitude'),
       ('engineLoad', 'engineLoad'),
       ('intakeTemp', 'intakeTemp'),
       ('coolantTemp', 'coolantTemp'),
       ('engineSpeed', 'engineSpeed'),
       ('batteryVoltage', 'batteryVoltage'),
       ('intakePressure', 'intakePressure'),
       ('throttlePosition', 'throttlePosition'),
       ('longTermFuelTrim1', 'longTermFuelTrim1'),
       ('shortTermFuelTrim1', 'shortTermFuelTrim1'),
       ('fuelPercentRemaining', 'fuelPercentRemaining'),
       ('acceleratorPedalPositionD', 'acceleratorPedalPositionD'),
       ('acceleratorPedalPositionE', 'acceleratorPedalPositionE'),
       ('dtc', 'dtc'),
       ('odometer', 'odometer'),
       ('barometricPressure', 'barometricPressure'),
       ('soc', 'soc'),
       ('chargeLimit', 'chargeLimit'),
       ('charging', 'charging'),
       ('charger', 'charger'),
       ('tires', 'tires'),
       ('oil', 'oil'),
       ('ambientTemp', 'ambientTemp'),
       ('range', 'range');


CREATE TABLE IF NOT EXISTS vehicle_signals_tracking_events_properties
(
    integration_id char(27),
    device_make_id character(27) ,
    property_id char(27) null,
    model character varying(100) NOT NULL,
    year int NOT NULL,
    description varchar(500) null,
    count int not null,
    created_at     timestamptz not null default current_timestamp,
    primary key (integration_id, device_make_id, property_id)
    );

CREATE TABLE IF NOT EXISTS vehicle_signals_tracking_events_missing_properties
(
    integration_id char(27),
    device_make_id character(27) ,
    property_id char(27) null,
    model character varying(100) NOT NULL,
    year int NOT NULL,
    description varchar(500) null,
    count int not null,
    created_at     timestamptz not null default current_timestamp,
    primary key (integration_id, device_make_id, property_id)
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

SET search_path = device_data_api, public;
DROP TABLE vehicle_signals_tracking_properties;
DROP TABLE vehicle_signals_tracking_events_properties;
DROP TABLE vehicle_signals_tracking_events_missing_properties;
-- +goose StatementEnd

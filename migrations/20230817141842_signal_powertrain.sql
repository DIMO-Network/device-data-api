-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path = device_data_api, public;
alter table vehicle_signals_available_properties
    alter column id type varchar(50) using id::varchar(50);

UPDATE vehicle_signals_available_properties SET power_train_type = '{ ICE, HEV, PHEV }' WHERE id LIKE 'maf';
UPDATE vehicle_signals_available_properties SET power_train_type = '{ ICE, HEV, PHEV }' WHERE id LIKE 'runTime';
UPDATE vehicle_signals_available_properties SET power_train_type = '{ ICE, HEV, PHEV }' WHERE id LIKE 'engineLoad';
UPDATE vehicle_signals_available_properties SET power_train_type = '{ ICE, HEV, PHEV }' WHERE id LIKE 'intakeTemp';
UPDATE vehicle_signals_available_properties SET power_train_type = '{ ICE, HEV, PHEV }' WHERE id LIKE 'coolantTemp';
UPDATE vehicle_signals_available_properties SET power_train_type = '{ ICE, HEV, PHEV }' WHERE id LIKE 'engineSpeed';
UPDATE vehicle_signals_available_properties SET power_train_type = '{ ICE, HEV, PHEV }' WHERE id LIKE 'intakePressure';
UPDATE vehicle_signals_available_properties SET power_train_type = '{ ICE, HEV, PHEV }' WHERE id LIKE 'longTermFuelTrim1';
UPDATE vehicle_signals_available_properties SET power_train_type = '{ ICE, HEV, PHEV }' WHERE id LIKE 'shortTermFuelTrim1';
UPDATE vehicle_signals_available_properties SET power_train_type = '{ ICE, HEV, PHEV }' WHERE id LIKE 'fuelPercentRemaining';
UPDATE vehicle_signals_available_properties SET power_train_type = '{ ICE, HEV, PHEV }' WHERE id LIKE 'oil';
UPDATE vehicle_signals_available_properties SET power_train_type = '{ BEV, PHEV }' WHERE id LIKE 'soc';
UPDATE vehicle_signals_available_properties SET power_train_type = '{ BEV, PHEV }' WHERE id LIKE 'chargeLimit';
UPDATE vehicle_signals_available_properties SET power_train_type = '{ BEV, PHEV }' WHERE id LIKE 'charging';


INSERT INTO vehicle_signals_available_properties(Id, name, power_train_type)
VALUES ('hybridEVBatteryCapacity', 'hybridEVBatteryCapacity', '{ HEV, PHEV }'),
       ('hybridEVBatteryChargeVoltage', 'hybridEVBatteryChargeVoltage', '{ HEV, PHEV }'),
       ('hybridEVBatteryMaxCellVolrage', 'hybridEVBatteryMaxCellVolrage', '{ HEV, PHEV }'),
       ('hybridEVBatteryRemainingChargeTime', 'hybridEVBatteryRemainingChargeTime', '{ HEV, PHEV }'),
       ('hybridEVBatterySoh', 'hybridEVBatterySoh', '{ HEV, PHEV }'),
       ('hybridEVBatteryRemaining', 'hybridEVBatteryRemaining', '{ HEV, PHEV }'),
       ('hybridEVBatteryTemp', 'hybridEVBatteryTemp', '{ HEV, PHEV }'),
       ('hybridEVEngineRPM', 'hybridEVEngineRPM', '{ HEV, PHEV }');


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
SET search_path = device_data_api, public;

-- +goose StatementEnd

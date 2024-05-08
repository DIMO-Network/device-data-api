-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path = device_data_api, public;
alter table vehicle_signals_available_properties add column valid_min_length int;

update vehicle_signals_available_properties set valid_min_length = 17 where name = 'vin';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
SET search_path = device_data_api, public;
alter table vehicle_signals_available_properties drop column valid_min_length;
-- +goose StatementEnd
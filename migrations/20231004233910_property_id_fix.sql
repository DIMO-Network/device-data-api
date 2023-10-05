-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path = device_data_api, public;
alter table report_vehicle_signals_events_tracking alter column property_id type varchar(50);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
SET search_path = device_data_api, public;
alter table report_vehicle_signals_events_tracking alter column property_id type char(27);
-- +goose StatementEnd
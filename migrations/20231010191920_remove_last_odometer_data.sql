-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path = device_data_api, public;
alter table user_device_data drop column last_odometer_event_at;
alter table user_device_data drop column real_last_odometer_event_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
SET search_path = device_data_api, public;
alter table user_device_data add column last_odometer_event_at timestamp with time zone;
alter table user_device_data add column real_last_odometer_event_at timestamp with time zone;
-- +goose StatementEnd
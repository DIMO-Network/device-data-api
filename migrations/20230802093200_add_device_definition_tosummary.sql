-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path = device_data_api, public;
ALTER TABLE report_vehicle_signals_events_summary
    ADD COLUMN device_definition_count integer not null DEFAULT 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE report_vehicle_signals_events_summary DROP COLUMN device_definition_count
-- +goose StatementEnd

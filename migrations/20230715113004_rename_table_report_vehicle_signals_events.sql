-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path = device_data_api, public;

ALTER TABLE report_vehicle_signals_events_properties RENAME TO report_vehicle_signals_events_tracking;
ALTER TABLE report_vehicle_signals_events RENAME TO report_vehicle_signals_events_all;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

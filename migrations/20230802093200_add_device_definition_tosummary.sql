-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path = device_data_api, public;
ALTER TABLE report_vehicle_signals_events_summary
    ADD COLUMN device_definition_id character(27) DEFAULT '';

ALTER TABLE report_vehicle_signals_events_summary DROP CONSTRAINT report_vehicle_signals_events_summary_pkey;
ALTER TABLE report_vehicle_signals_events_summary ADD CONSTRAINT report_vehicle_signals_events_summary_pkey PRIMARY KEY (date_id, integration_id, power_train_type, device_definition_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE report_vehicle_signals_events_summary DROP COLUMN device_definition_id
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path = device_data_api, public;
alter table report_vehicle_signals_events_tracking
    alter column date_id type char(8) using date_id::char(8);

alter table report_vehicle_signals_events_all
    alter column date_id type char(8) using date_id::char(8);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
SET search_path = device_data_api, public;

alter table report_vehicle_signals_events_tracking
    alter column date_id type char(27) using date_id::char(27);

alter table report_vehicle_signals_events_all
    alter column date_id type char(27) using date_id::char(27);
-- +goose StatementEnd

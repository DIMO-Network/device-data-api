-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS  report_vehicle_signals_events_user_devices (
    date_id character(27) not null,
    integration_id character(27) not null,
    power_train_type character(4) not null,
    count integer not null,
    created_at timestamp with time zone not null default CURRENT_TIMESTAMP,
    primary key (date_id, integration_id, power_train_type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table report_vehicle_signals_events_user_devices;
-- +goose StatementEnd

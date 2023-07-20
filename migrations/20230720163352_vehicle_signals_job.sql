-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path = device_data_api, public;
CREATE TABLE IF NOT EXISTS vehicle_signals_job
(
    vehicle_signals_job_id char(27),
    start_date timestamptz not null,
    end_date timestamptz not null,
    created_at     timestamptz not null default current_timestamp,
    primary key (vehicle_signals_job_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table vehicle_signals_job;
-- +goose StatementEnd

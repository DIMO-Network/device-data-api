-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

SET search_path = device_data_api, public;

CREATE TABLE IF NOT EXISTS user_device_data
(
    user_device_id char(27), -- ksuid
    error_data           jsonb,
    last_odometer_event_at           timestamptz,
    integration_id char(27), -- ksuid
    real_last_odometer_event_at timestamptz,
    signals           jsonb,
    created_at     timestamptz not null default current_timestamp,
    updated_at     timestamptz not null default current_timestamp,
    primary key (user_device_id, integration_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

SET search_path = device_data_api, public;
DROP TABLE user_device_data;

-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

SET search_path = devices_data_api, public;

CREATE TABLE IF NOT EXISTS user_device_data
(
    user_device_id char(27) PRIMARY KEY, -- ksuid
    data           jsonb,
    error_data           jsonb,
    last_odometer_event_at           timestamptz,
    integration_id char(27),
    real_last_odometer_event_at timestamptz,
    signals           jsonb,
    created_at     timestamptz not null default current_timestamp,
    updated_at     timestamptz not null default current_timestamp
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

SET search_path = devices_data_api, public;
DROP TABLE user_device_data;

-- +goose StatementEnd

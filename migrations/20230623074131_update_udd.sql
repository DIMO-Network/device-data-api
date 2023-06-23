-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path = device_data_api, public;

alter table user_device_data alter column integration_id set not null;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
SET search_path = device_data_api, public;

alter table user_device_data alter column integration_id drop not null;
-- +goose StatementEnd

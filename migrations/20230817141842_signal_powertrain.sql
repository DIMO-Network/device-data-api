-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path = device_data_api, public;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
SET search_path = device_data_api, public;

-- +goose StatementEnd

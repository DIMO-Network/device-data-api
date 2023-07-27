-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SET search_path = device_data_api, public;
ALTER TABLE vehicle_signals_available_properties ADD COLUMN power_train_type TEXT[] DEFAULT '{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE vehicle_signals_available_properties DROP COLUMN power_train_type;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE subscriptions ALTER COLUMN end_date DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

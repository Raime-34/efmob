-- +goose Up
-- +goose StatementBegin
ALTER TABLE subscriptions
ADD COLUMN price INT NOT NULL,
ADD CONSTRAINT subscriptions_price_positive CHECK (price > 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

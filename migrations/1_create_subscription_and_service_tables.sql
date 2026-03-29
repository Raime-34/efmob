-- +goose Up
-- +goose StatementBegin
CREATE TABLE services (
    service_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE subscriptions (
    user_id UUID NOT NULL,
    service_id INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    PRIMARY KEY (user_id, service_id),
    FOREIGN KEY (service_id) REFERENCES services(service_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

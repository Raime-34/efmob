-- +goose Up
-- +goose StatementBegin
ALTER TABLE services
ADD CONSTRAINT uq_name UNIQUE (name);

CREATE OR REPLACE FUNCTION get_or_create_service(service_name VARCHAR(255)) RETURNS INT AS
$$
DECLARE
    result_id INT;
BEGIN
    INSERT INTO services (name) VALUES (service_name)
    ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
    RETURNING service_id INTO result_id;
    RETURN result_id;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS get_or_create_service(VARCHAR);
-- +goose StatementEnd

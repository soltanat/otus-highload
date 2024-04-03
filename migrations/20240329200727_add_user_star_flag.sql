-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN star BOOLEAN DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN star;
-- +goose StatementEnd

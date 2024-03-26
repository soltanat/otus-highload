-- +goose Up
-- +goose StatementBegin
create index ix_users_name on users (second_name collate "C", first_name collate "C");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index ix_users_name;
-- +goose StatementEnd

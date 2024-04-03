-- +goose Up
-- +goose StatementBegin
create table if not exists friends
(
    user_id   uuid references users (id),
    friend_id uuid references users (id),
    star      boolean default false,
    primary key (user_id, friend_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists friends;
-- +goose StatementEnd

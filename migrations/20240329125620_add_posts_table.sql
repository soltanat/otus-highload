-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts
(
    id         uuid PRIMARY KEY,
    author_id  uuid                     NOT NULL references users (id),
    created_at timestamp with time zone NOT NULL,
    text    TEXT                     NOT NULL
);

CREATE INDEX posts_author_id_idx ON posts (author_id, created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
-- +goose StatementEnd

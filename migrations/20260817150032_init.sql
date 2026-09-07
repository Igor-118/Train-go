-- +goose Up
CREATE TABLE users (
    id bigint primary key,
    city text,
    created_at timestamp default NOW()
);

-- +goose Down


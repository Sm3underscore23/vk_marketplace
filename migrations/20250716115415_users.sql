-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    user_login    VARCHAR(32) UNIQUE NOT NULL,
    password_hash VARCHAR(64) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
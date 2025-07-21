-- +goose Up
-- +goose StatementBegin
CREATE TABLE advertisements (
    id SERIAL PRIMARY KEY,
    title VARCHAR(32) NOT NULL,
    description TEXT NOT NULL,
    image_url VARCHAR(256) NOT NULL,
    price INT NOT NULL,
    author_id INT NOT NULL,
    FOREIGN KEY(author_id) REFERENCES users(id),
    created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE advertisements;
-- +goose StatementEnd
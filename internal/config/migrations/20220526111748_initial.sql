-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls (
    id TEXT UNIQUE NOT NULL PRIMARY KEY,
    short_url TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    user_id UUID NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP DATABASE urls;
-- +goose StatementEnd

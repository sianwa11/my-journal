-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN bio TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN bio;
-- +goose StatementEnd

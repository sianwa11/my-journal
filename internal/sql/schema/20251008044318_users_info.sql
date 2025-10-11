-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN email TEXT DEFAULT '';
ALTER TABLE users ADD COLUMN github TEXT DEFAULT '';
ALTER TABLE users ADD COLUMN linkedin TEXT DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN email;
ALTER TABLE users DROP COLUMN github;
ALTER TABLE users DROP COLUMN linkedin;
-- +goose StatementEnd
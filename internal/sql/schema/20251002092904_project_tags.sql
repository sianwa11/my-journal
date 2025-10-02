-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_tags(
  project_id INTEGER NOT NULL,
  tag_id INTEGER NOT NULL,
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (project_id, tag_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP IF EXISTS TABLE project_tags;
-- +goose StatementEnd

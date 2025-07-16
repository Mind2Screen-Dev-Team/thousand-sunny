-- +goose Up
-- +goose StatementBegin
CREATE TABLE example_users(
  id UUID PRIMARY KEY,
  name VARCHAR NOT NULL,
  level INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX on example_users(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE example_users;
-- +goose StatementEnd

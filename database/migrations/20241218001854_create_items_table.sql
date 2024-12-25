-- +goose Up
-- +goose StatementBegin
CREATE TABLE items (
  id UUID PRIMARY KEY,
  name VARCHAR NOT NULL,
  value INT NOT NULL
);

CREATE UNIQUE INDEX ON items(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
-- +goose StatementEnd

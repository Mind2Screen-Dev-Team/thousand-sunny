-- +goose Up
-- +goose StatementBegin
CREATE TABLE players(
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  level INT NOT NULL,
  class VARCHAR NOT NULL,
  gold bigint NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX on players(name);
ALTER TABLE players ADD CONSTRAINT chk_gold_non_negative CHECK (gold >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE players;
-- +goose StatementEnd

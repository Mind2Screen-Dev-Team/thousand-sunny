-- +goose Up
-- +goose StatementBegin
CREATE TABLE inventories (
  player_id BIGINT NOT NULL REFERENCES players(id),
  item_id UUID NOT NULL REFERENCES items(id),
  PRIMARY KEY(player_id, item_id)
);

CREATE INDEX ON inventories(player_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE inventories;
-- +goose StatementEnd

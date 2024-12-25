-- +goose Up
-- +goose StatementBegin
INSERT INTO items (id, name, value) VALUES
('42c6294c-56de-49d2-be2e-055b2a2151a6', 'Rusty Sword', 100),
('2e9e9593-c5ec-4554-9e15-131aa0b63127', 'Bat Wing',6),
('0faffedc-2047-4616-9357-51d22fe80ff7', 'Skeleton Femur', 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM items WHERE id IN ('42c6294c-56de-49d2-be2e-055b2a2151a6','2e9e9593-c5ec-4554-9e15-131aa0b63127','0faffedc-2047-4616-9357-51d22fe80ff7');
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
INSERT INTO example_users (id, name, level, created_at, updated_at) VALUES
('0198121c-a1db-79a9-bc37-44abd13ff402', 'Sephiro', 32, '2024-01-15 14:23:45Z', '2024-02-10 08:30:15Z'),
('0198121c-d011-73c1-a578-7025415cc3c4', 'Spell Queen', 47, '2023-12-05 10:15:12Z', '2024-01-22 16:45:22Z'),
('0198121c-e829-740d-9d8c-04493e96ba8a', 'Storm Alex', 29, '2023-11-20 09:20:34Z', '2024-01-10 12:10:58Z'),
('0198121c-ff30-7bf4-98cc-1079c62345dd', 'Shadow Master', 50, '2023-10-08 19:50:14Z', '2023-12-18 11:35:45Z'),
('0198121d-160e-7763-9203-0f82ec372efa', 'Fire Furry', 55, '2023-09-03 17:45:23Z', '2023-11-14 14:20:30Z');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM example_users WHERE id IN ('0198121c-a1db-79a9-bc37-44abd13ff402','0198121c-d011-73c1-a578-7025415cc3c4','0198121c-e829-740d-9d8c-04493e96ba8a','0198121c-ff30-7bf4-98cc-1079c62345dd','0198121d-160e-7763-9203-0f82ec372efa');
-- +goose StatementEnd

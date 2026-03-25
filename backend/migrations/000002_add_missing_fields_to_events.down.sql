-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
    DROP COLUMN IF EXISTS rules,
    DROP COLUMN IF EXISTS recommendations,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS max_participants;
-- +goose StatementEnd
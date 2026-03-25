-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
    ADD COLUMN IF NOT EXISTS rules TEXT,
    ADD COLUMN IF NOT EXISTS recommendations TEXT,
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'finished', 'cancelled')),
    ADD COLUMN IF NOT EXISTS max_participants INTEGER DEFAULT 0;

COMMENT ON COLUMN events.rules IS 'Правила игры';
COMMENT ON COLUMN events.recommendations IS 'Рекомендации по подаркам';
COMMENT ON COLUMN events.status IS 'Статус события (draft / active / finished / cancelled)';
COMMENT ON COLUMN events.max_participants IS 'Максимальное количество участников (0 = без ограничения)';
-- +goose StatementEnd
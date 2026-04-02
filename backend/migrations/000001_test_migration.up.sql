-- +goose Up
-- +goose StatementBegin

-- ================================================================
-- FIXED & ALIGNED by Grok — 02.04.2026
-- Полная миграция с COMMENT ON COLUMN для КАЖДОЙ колонки
-- 13 колонок в events + подробные описания
-- ================================================================

-- Расширение для UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ==================== USERS ====================
CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    email           TEXT UNIQUE NOT NULL,
    oauth_id        TEXT UNIQUE NOT NULL,
    oauth_provider  TEXT NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE users IS 'Пользователи системы';
COMMENT ON COLUMN users.id IS 'Уникальный идентификатор пользователя';
COMMENT ON COLUMN users.name IS 'Отображаемое имя пользователя';
COMMENT ON COLUMN users.email IS 'Email пользователя';
COMMENT ON COLUMN users.oauth_id IS 'Идентификатор пользователя у провайдера OAuth';
COMMENT ON COLUMN users.oauth_provider IS 'Название провайдера (github, vk, google и т.д.)';
COMMENT ON COLUMN users.created_at IS 'Дата создания записи';
COMMENT ON COLUMN users.updated_at IS 'Дата последнего обновления записи';

-- ==================== EVENTS ====================
CREATE TABLE IF NOT EXISTS events (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title               TEXT NOT NULL,
    description         TEXT,
    rules               TEXT,
    recommendations     TEXT,
    organizer_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date          TIMESTAMPTZ NOT NULL,
    draw_date           TIMESTAMPTZ,
    end_date            TIMESTAMPTZ NOT NULL,
    status              TEXT NOT NULL DEFAULT 'draft',
    max_participants    INT NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE events IS 'События "Тайный Санта"';
COMMENT ON COLUMN events.id IS 'Уникальный идентификатор события';
COMMENT ON COLUMN events.title IS 'Название события';
COMMENT ON COLUMN events.description IS 'Описание события';
COMMENT ON COLUMN events.rules IS 'Правила события';
COMMENT ON COLUMN events.recommendations IS 'Рекомендации и пожелания по подаркам';
COMMENT ON COLUMN events.organizer_id IS 'ID организатора события';
COMMENT ON COLUMN events.start_date IS 'Дата начала события';
COMMENT ON COLUMN events.draw_date IS 'Дата проведения жеребьёвки';
COMMENT ON COLUMN events.end_date IS 'Дата окончания события';
COMMENT ON COLUMN events.status IS 'Текущий статус события (draft, invitation_open и т.д.)';
COMMENT ON COLUMN events.max_participants IS 'Максимальное количество участников';
COMMENT ON COLUMN events.created_at IS 'Дата создания события';
COMMENT ON COLUMN events.updated_at IS 'Дата последнего обновления события';

-- ==================== PARTICIPANTS ====================
CREATE TABLE IF NOT EXISTS participants (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id      UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role          TEXT NOT NULL DEFAULT 'participant',
    gift_sent     BOOLEAN NOT NULL DEFAULT false,
    gift_sent_at  TIMESTAMPTZ,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (event_id, user_id)
);

COMMENT ON TABLE participants IS 'Участники событий';
COMMENT ON COLUMN participants.id IS 'Уникальный идентификатор участника';
COMMENT ON COLUMN participants.event_id IS 'Ссылка на событие';
COMMENT ON COLUMN participants.user_id IS 'Ссылка на пользователя';
COMMENT ON COLUMN participants.role IS 'Роль участника (organizer или participant)';
COMMENT ON COLUMN participants.gift_sent IS 'Подарок уже отправлен?';
COMMENT ON COLUMN participants.gift_sent_at IS 'Когда отметили отправку подарка';
COMMENT ON COLUMN participants.created_at IS 'Дата добавления участника';
COMMENT ON COLUMN participants.updated_at IS 'Дата последнего обновления участника';

-- ==================== ASSIGNMENTS ====================
CREATE TABLE IF NOT EXISTS assignments (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id     UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    giver_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (event_id, giver_id)
);

COMMENT ON TABLE assignments IS 'Результат жеребьёвки (кто кому дарит)';
COMMENT ON COLUMN assignments.id IS 'Уникальный ID назначения';
COMMENT ON COLUMN assignments.event_id IS 'Событие';
COMMENT ON COLUMN assignments.giver_id IS 'Кто дарит (Санта)';
COMMENT ON COLUMN assignments.receiver_id IS 'Кому дарит';
COMMENT ON COLUMN assignments.created_at IS 'Дата создания назначения';

-- ==================== WISHLISTS ====================
CREATE TABLE IF NOT EXISTS wishlists (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id  UUID NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    visibility      TEXT NOT NULL DEFAULT 'santa_only',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE wishlists IS 'Вишлисты участников';
COMMENT ON COLUMN wishlists.id IS 'Уникальный идентификатор вишлиста';
COMMENT ON COLUMN wishlists.participant_id IS 'Участник, которому принадлежит вишлист';
COMMENT ON COLUMN wishlists.visibility IS 'Видимость (public, friends, santa_only)';
COMMENT ON COLUMN wishlists.created_at IS 'Дата создания вишлиста';
COMMENT ON COLUMN wishlists.updated_at IS 'Дата последнего обновления вишлиста';

-- ==================== WISHLIST ITEMS ====================
CREATE TABLE IF NOT EXISTS wishlist_items (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wishlist_id  UUID NOT NULL REFERENCES wishlists(id) ON DELETE CASCADE,
    title        TEXT NOT NULL,
    link         TEXT,
    image_url    TEXT,
    comment      TEXT,
    created_at   TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE wishlist_items IS 'Элементы вишлиста';
COMMENT ON COLUMN wishlist_items.id IS 'Уникальный идентификатор элемента';
COMMENT ON COLUMN wishlist_items.wishlist_id IS 'Ссылка на вишлист';
COMMENT ON COLUMN wishlist_items.title IS 'Название желаемого подарка';
COMMENT ON COLUMN wishlist_items.link IS 'Ссылка на подарок';
COMMENT ON COLUMN wishlist_items.image_url IS 'Ссылка на изображение подарка';
COMMENT ON COLUMN wishlist_items.comment IS 'Комментарий пользователя';
COMMENT ON COLUMN wishlist_items.created_at IS 'Дата добавления товара';

-- ==================== INVITATIONS ====================
CREATE TABLE IF NOT EXISTS invitations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id    UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    token       TEXT UNIQUE NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_by  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE invitations IS 'Приглашения по ссылке (многоразовые)';
COMMENT ON COLUMN invitations.id IS 'Уникальный идентификатор приглашения';
COMMENT ON COLUMN invitations.event_id IS 'Событие';
COMMENT ON COLUMN invitations.token IS 'Уникальный токен для ссылки';
COMMENT ON COLUMN invitations.expires_at IS 'Срок действия приглашения';
COMMENT ON COLUMN invitations.created_by IS 'Организатор, который создал приглашение';
COMMENT ON COLUMN invitations.created_at IS 'Дата создания';
COMMENT ON COLUMN invitations.updated_at IS 'Дата последнего обновления';

-- ==================== MESSAGES ====================
CREATE TABLE IF NOT EXISTS messages (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id    UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    sender_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content     TEXT NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE messages IS 'Сообщения в чате между Сантой и получателем';
COMMENT ON COLUMN messages.id IS 'Уникальный ID сообщения';
COMMENT ON COLUMN messages.event_id IS 'Событие';
COMMENT ON COLUMN messages.sender_id IS 'Отправитель';
COMMENT ON COLUMN messages.receiver_id IS 'Получатель';
COMMENT ON COLUMN messages.content IS 'Текст сообщения';
COMMENT ON COLUMN messages.created_at IS 'Дата отправки';

-- Индексы
CREATE INDEX idx_users_oauth ON users(oauth_id, oauth_provider);
CREATE INDEX idx_events_organizer ON events(organizer_id);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_participants_event ON participants(event_id);
CREATE INDEX idx_assignments_event ON assignments(event_id);
CREATE INDEX idx_wishlists_participant ON wishlists(participant_id);
CREATE INDEX idx_messages_event_pair ON messages(event_id, sender_id, receiver_id);

-- +goose StatementEnd
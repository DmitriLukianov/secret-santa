-- +goose Up
-- +goose StatementBegin

-- Включаем расширение для UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ==================== USERS ====================
CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    email           TEXT UNIQUE NOT NULL,
    oauth_id        TEXT UNIQUE NOT NULL DEFAULT '',
    oauth_provider  TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now()
);

COMMENT ON TABLE users IS 'Пользователи системы';
COMMENT ON COLUMN users.id IS 'Уникальный идентификатор пользователя';
COMMENT ON COLUMN users.name IS 'Имя пользователя';
COMMENT ON COLUMN users.email IS 'Email пользователя';
COMMENT ON COLUMN users.oauth_id IS 'Идентификатор пользователя у провайдера OAuth';
COMMENT ON COLUMN users.oauth_provider IS 'Название провайдера OAuth';
COMMENT ON COLUMN users.created_at IS 'Дата создания';
COMMENT ON COLUMN users.updated_at IS 'Дата последнего обновления';

CREATE INDEX IF NOT EXISTS idx_users_oauth ON users(oauth_id, oauth_provider);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- ==================== EVENTS ====================
CREATE TABLE IF NOT EXISTS events (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title               TEXT NOT NULL,
    description         TEXT,
    rules               TEXT,
    recommendations     TEXT,
    organizer_id        UUID NOT NULL,
    start_date          TIMESTAMPTZ NOT NULL,
    draw_date           TIMESTAMPTZ,
    end_date            TIMESTAMPTZ NOT NULL,
    status              TEXT NOT NULL DEFAULT 'draft',
    max_participants    INT NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ DEFAULT now(),
    updated_at          TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_events_organizer 
        FOREIGN KEY (organizer_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);

COMMENT ON TABLE events IS 'События "Тайный Санта"';
COMMENT ON COLUMN events.id IS 'Уникальный идентификатор события';
COMMENT ON COLUMN events.title IS 'Название события';
COMMENT ON COLUMN events.description IS 'Описание события';
COMMENT ON COLUMN events.rules IS 'Правила события';
COMMENT ON COLUMN events.recommendations IS 'Рекомендации и пожелания';
COMMENT ON COLUMN events.organizer_id IS 'Организатор события';
COMMENT ON COLUMN events.start_date IS 'Дата начала события';
COMMENT ON COLUMN events.draw_date IS 'Дата проведения жеребьёвки';
COMMENT ON COLUMN events.end_date IS 'Дата окончания события';
COMMENT ON COLUMN events.status IS 'Статус события (draft, active, finished, cancelled)';
COMMENT ON COLUMN events.max_participants IS 'Максимальное количество участников';
COMMENT ON COLUMN events.created_at IS 'Дата создания';
COMMENT ON COLUMN events.updated_at IS 'Дата последнего обновления';

CREATE INDEX IF NOT EXISTS idx_events_organizer ON events(organizer_id);
CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);

-- ==================== PARTICIPANTS ====================
CREATE TABLE IF NOT EXISTS participants (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id      UUID NOT NULL,
    user_id       UUID NOT NULL,
    role          TEXT NOT NULL DEFAULT 'participant',
    gift_sent     BOOLEAN NOT NULL DEFAULT false,
    gift_sent_at  TIMESTAMPTZ,
    created_at    TIMESTAMPTZ DEFAULT now(),
    updated_at    TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_participants_event 
        FOREIGN KEY (event_id) 
        REFERENCES events(id) 
        ON DELETE CASCADE,

    CONSTRAINT fk_participants_user 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,

    CONSTRAINT unique_event_user UNIQUE (event_id, user_id)
);

COMMENT ON TABLE participants IS 'Участники событий';
COMMENT ON COLUMN participants.id IS 'Уникальный идентификатор участника';
COMMENT ON COLUMN participants.event_id IS 'Событие';
COMMENT ON COLUMN participants.user_id IS 'Пользователь';
COMMENT ON COLUMN participants.role IS 'Роль участника (organizer / participant)';
COMMENT ON COLUMN participants.gift_sent IS 'Подарок отправлен';
COMMENT ON COLUMN participants.gift_sent_at IS 'Дата отправки подарка';
COMMENT ON COLUMN participants.created_at IS 'Дата создания';
COMMENT ON COLUMN participants.updated_at IS 'Дата последнего обновления';

CREATE INDEX IF NOT EXISTS idx_participants_event ON participants(event_id);
CREATE INDEX IF NOT EXISTS idx_participants_user ON participants(user_id);
CREATE INDEX IF NOT EXISTS idx_participants_role ON participants(role);

-- ==================== ASSIGNMENTS ====================
CREATE TABLE IF NOT EXISTS assignments (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id     UUID NOT NULL,
    giver_id     UUID NOT NULL,
    receiver_id  UUID NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_assignments_event 
        FOREIGN KEY (event_id) 
        REFERENCES events(id) 
        ON DELETE CASCADE,

    CONSTRAINT fk_assignments_giver 
        FOREIGN KEY (giver_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,

    CONSTRAINT fk_assignments_receiver 
        FOREIGN KEY (receiver_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,

    CONSTRAINT unique_event_giver UNIQUE (event_id, giver_id)
);

COMMENT ON TABLE assignments IS 'Назначения "кто кому дарит" (жеребьёвка)';
COMMENT ON COLUMN assignments.id IS 'Уникальный идентификатор назначения';
COMMENT ON COLUMN assignments.event_id IS 'Событие';
COMMENT ON COLUMN assignments.giver_id IS 'Кто дарит';
COMMENT ON COLUMN assignments.receiver_id IS 'Кому дарит';
COMMENT ON COLUMN assignments.created_at IS 'Дата создания';

CREATE INDEX IF NOT EXISTS idx_assignments_event ON assignments(event_id);
CREATE INDEX IF NOT EXISTS idx_assignments_giver ON assignments(giver_id);
CREATE INDEX IF NOT EXISTS idx_assignments_receiver ON assignments(receiver_id);

-- ==================== WISHLISTS ====================
CREATE TABLE IF NOT EXISTS wishlists (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id  UUID NOT NULL,
    visibility      TEXT NOT NULL DEFAULT 'santa_only',
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_wishlists_participant 
        FOREIGN KEY (participant_id) 
        REFERENCES participants(id) 
        ON DELETE CASCADE
);

COMMENT ON TABLE wishlists IS 'Вишлисты участников';
COMMENT ON COLUMN wishlists.id IS 'Уникальный идентификатор вишлиста';
COMMENT ON COLUMN wishlists.participant_id IS 'Участник';
COMMENT ON COLUMN wishlists.visibility IS 'Видимость вишлиста (all / friends / santa_only)';
COMMENT ON COLUMN wishlists.created_at IS 'Дата создания';
COMMENT ON COLUMN wishlists.updated_at IS 'Дата последнего обновления';

CREATE INDEX IF NOT EXISTS idx_wishlists_participant ON wishlists(participant_id);

-- ==================== WISHLIST ITEMS ====================
CREATE TABLE IF NOT EXISTS wishlist_items (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wishlist_id  UUID NOT NULL,
    title        TEXT NOT NULL,
    link         TEXT,
    image_url    TEXT,
    comment      TEXT,
    created_at   TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_wishlist_items_wishlist 
        FOREIGN KEY (wishlist_id) 
        REFERENCES wishlists(id) 
        ON DELETE CASCADE
);

COMMENT ON TABLE wishlist_items IS 'Элементы вишлиста';
COMMENT ON COLUMN wishlist_items.id IS 'Уникальный идентификатор элемента';
COMMENT ON COLUMN wishlist_items.wishlist_id IS 'Вишлист';
COMMENT ON COLUMN wishlist_items.title IS 'Название желаемого подарка';
COMMENT ON COLUMN wishlist_items.link IS 'Ссылка на подарок';
COMMENT ON COLUMN wishlist_items.image_url IS 'Ссылка на изображение';
COMMENT ON COLUMN wishlist_items.comment IS 'Комментарий';
COMMENT ON COLUMN wishlist_items.created_at IS 'Дата добавления';

CREATE INDEX IF NOT EXISTS idx_wishlist_items_wishlist ON wishlist_items(wishlist_id);

CREATE TABLE IF NOT EXISTS invitations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id    UUID NOT NULL,
    token       TEXT UNIQUE NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_by  UUID NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_invitations_event 
        FOREIGN KEY (event_id) 
        REFERENCES events(id) 
        ON DELETE CASCADE,

    CONSTRAINT fk_invitations_created_by 
        FOREIGN KEY (created_by) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);

COMMENT ON TABLE invitations IS 'Приглашения по ссылке для событий Тайный Санта (многоразовые)';
COMMENT ON COLUMN invitations.id IS 'Уникальный идентификатор приглашения';
COMMENT ON COLUMN invitations.event_id IS 'Событие, к которому приглашение';
COMMENT ON COLUMN invitations.token IS 'Уникальный токен для ссылки';
COMMENT ON COLUMN invitations.expires_at IS 'Срок действия приглашения';
COMMENT ON COLUMN invitations.created_by IS 'Кто создал приглашение (организатор)';
COMMENT ON COLUMN invitations.created_at IS 'Дата создания';
COMMENT ON COLUMN invitations.updated_at IS 'Дата последнего обновления';

CREATE INDEX IF NOT EXISTS idx_invitations_token ON invitations(token);
CREATE INDEX IF NOT EXISTS idx_invitations_event ON invitations(event_id);

-- +goose StatementEnd
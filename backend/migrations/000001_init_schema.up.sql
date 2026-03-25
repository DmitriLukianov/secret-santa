-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ==================== USERS ====================
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    email           TEXT UNIQUE NOT NULL,
    oauth_id        TEXT UNIQUE NOT NULL DEFAULT '',
    oauth_provider  TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_oauth ON users(oauth_id, oauth_provider);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- ==================== EVENTS ====================
-- ==================== EVENTS ====================
CREATE TABLE events (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title               TEXT NOT NULL,                    -- было name → теперь title
    description         TEXT,
    rules               TEXT,
    recommendations     TEXT,
    organizer_id        UUID NOT NULL,
    start_date          TIMESTAMPTZ NOT NULL,             -- теперь обязательно
    draw_date           TIMESTAMPTZ,
    end_date            TIMESTAMPTZ NOT NULL,             -- теперь обязательно
    status              TEXT NOT NULL DEFAULT 'draft',
    max_participants    INT NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ DEFAULT now(),
    updated_at          TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_events_organizer 
        FOREIGN KEY (organizer_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);

-- Индексы для быстрых поисков
CREATE INDEX IF NOT EXISTS idx_events_organizer ON events(organizer_id);
CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);
-- ==================== PARTICIPANTS ====================
CREATE TABLE participants (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id   UUID NOT NULL,
    user_id    UUID NOT NULL,
    role       TEXT NOT NULL DEFAULT 'participant',
    gift_sent  BOOLEAN NOT NULL DEFAULT false,
    gift_sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_participants_event FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    CONSTRAINT fk_participants_user  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_event_user UNIQUE (event_id, user_id)
);

-- ==================== ASSIGNMENTS ====================
CREATE TABLE assignments (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id     UUID NOT NULL,
    giver_id     UUID NOT NULL,
    receiver_id  UUID NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_assignments_event    FOREIGN KEY (event_id)    REFERENCES events(id)    ON DELETE CASCADE,
    CONSTRAINT fk_assignments_giver    FOREIGN KEY (giver_id)    REFERENCES users(id)    ON DELETE CASCADE,
    CONSTRAINT fk_assignments_receiver FOREIGN KEY (receiver_id) REFERENCES users(id)    ON DELETE CASCADE,
    CONSTRAINT unique_event_giver UNIQUE (event_id, giver_id)
);

-- ==================== WISHLISTS ====================
CREATE TABLE wishlists (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id UUID NOT NULL,           -- привязка к участнику события
    visibility    TEXT NOT NULL DEFAULT 'santa_only',
    created_at    TIMESTAMPTZ DEFAULT now(),
    updated_at    TIMESTAMPTZ DEFAULT now(),

    CONSTRAINT fk_wishlist_participant FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE
);

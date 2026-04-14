-- Пользователи системы
CREATE TABLE IF NOT EXISTS users (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(255)  NOT NULL,
    email          VARCHAR(255)  NOT NULL UNIQUE,
    oauth_id       VARCHAR(255)  NOT NULL,
    oauth_provider VARCHAR(50)   NOT NULL,
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_oauth ON users (oauth_id, oauth_provider);
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

-- События Тайный Санта
CREATE TABLE IF NOT EXISTS events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(255) NOT NULL,
    organizer_notes TEXT,
    organizer_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date      TIMESTAMPTZ  NOT NULL,
    draw_date       TIMESTAMPTZ,
    status          VARCHAR(50)  NOT NULL DEFAULT 'registration',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_events_organizer_id ON events (organizer_id);
CREATE INDEX IF NOT EXISTS idx_events_status ON events (status);

-- Участники события
CREATE TABLE IF NOT EXISTS participants (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id   UUID        NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (event_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_participants_event_id ON participants (event_id);
CREATE INDEX IF NOT EXISTS idx_participants_user_id ON participants (user_id);

-- Назначения жеребьёвки
CREATE TABLE IF NOT EXISTS assignments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id    UUID        NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    giver_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (event_id, giver_id),
    CONSTRAINT no_self_assignment CHECK (giver_id != receiver_id)
);

CREATE INDEX IF NOT EXISTS idx_assignments_event_id ON assignments (event_id);
CREATE INDEX IF NOT EXISTS idx_assignments_giver_id ON assignments (giver_id);

-- Вишлисты
CREATE TABLE IF NOT EXISTS wishlists (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id UUID        REFERENCES participants(id) ON DELETE CASCADE,
    user_id        UUID        REFERENCES users(id) ON DELETE CASCADE,
    visibility     VARCHAR(50) NOT NULL DEFAULT 'santa_only',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_wishlists_participant_id ON wishlists (participant_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_wishlists_user_id ON wishlists (user_id) WHERE user_id IS NOT NULL;

-- Элементы вишлиста
CREATE TABLE IF NOT EXISTS wishlist_items (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    wishlist_id UUID         NOT NULL REFERENCES wishlists(id) ON DELETE CASCADE,
    title       VARCHAR(500) NOT NULL,
    link        TEXT,
    image_url   TEXT,
    price       NUMERIC(10,2),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_wishlist_items_wishlist_id ON wishlist_items (wishlist_id);

-- Приглашения
CREATE TABLE IF NOT EXISTS invitations (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id   UUID         NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ  NOT NULL,
    created_by UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_invitations_event_id ON invitations (event_id);
CREATE INDEX IF NOT EXISTS idx_invitations_token ON invitations (token);

-- Сообщения анонимного чата
CREATE TABLE IF NOT EXISTS messages (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id    UUID        NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    sender_id   UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content     TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_messages_event_id ON messages (event_id);
CREATE INDEX IF NOT EXISTS idx_messages_pair ON messages (event_id, sender_id, receiver_id);

-- OTP-коды для входа по email
CREATE TABLE IF NOT EXISTS email_verification_codes (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email      VARCHAR(255) NOT NULL,
    code       VARCHAR(10)  NOT NULL,
    expires_at TIMESTAMPTZ  NOT NULL,
    used       BOOLEAN      NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_verification_email ON email_verification_codes (email);
CREATE INDEX IF NOT EXISTS idx_verification_lookup ON email_verification_codes (email, code, used, expires_at);

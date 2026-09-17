CREATE EXTENSION IF NOT EXISTS pgcrypto; -- provides gen_random_uuid()

CREATE TABLE players (
    player_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username    VARCHAR(32) NOT NULL UNIQUE,
    mu          DOUBLE PRECISION NOT NULL DEFAULT 25.0,
    sigma       DOUBLE PRECISION NOT NULL DEFAULT 8.333,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE matches (
    match_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status       VARCHAR(16) NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled')),
    map          VARCHAR(32) NOT NULL,
    started_at   TIMESTAMPTZ,
    ended_at     TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE match_participants (
    participant_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id          UUID NOT NULL REFERENCES matches(match_id),
    player_id         UUID NOT NULL REFERENCES players(player_id),
    team              CHAR(1) NOT NULL CHECK (team IN ('A', 'B')),
    result            VARCHAR(4) NOT NULL CHECK (result IN ('win', 'loss', 'draw')),
    character         VARCHAR(32) NOT NULL,
    kills             INT NOT NULL DEFAULT 0 CHECK (kills >= 0),
    deaths            INT NOT NULL DEFAULT 0 CHECK (deaths >= 0),
    assists           INT NOT NULL DEFAULT 0 CHECK (assists >= 0),
    performance_score DOUBLE PRECISION, -- nullable: computed by rating-service after the match
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT unique_player_per_match UNIQUE (match_id, player_id)
);

-- Explicit index needed for lookups filtering on player_id alone.
CREATE INDEX idx_participant_player_id ON match_participants(player_id);

CREATE TABLE rating_history (
    rating_history_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id     UUID NOT NULL UNIQUE REFERENCES match_participants(participant_id),
    mu_before          DOUBLE PRECISION NOT NULL,
    sigma_before       DOUBLE PRECISION NOT NULL,
    mu_after           DOUBLE PRECISION NOT NULL,
    sigma_after        DOUBLE PRECISION NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

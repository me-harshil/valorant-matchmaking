CREATE TABLE rank_tiers (
    tier_id     SERIAL PRIMARY KEY,
    min_rating  DOUBLE PRECISION NOT NULL,
    name        VARCHAR(16) NOT NULL,
    sub         SMALLINT NOT NULL CHECK (sub >= 0 AND sub <= 3),

    CONSTRAINT unique_name_sub UNIQUE (name, sub)
);

CREATE INDEX idx_rank_tiers_min_rating ON rank_tiers(min_rating);

INSERT INTO rank_tiers (min_rating, name, sub) VALUES
    (-100, 'Iron', 1), (5, 'Iron', 2), (10, 'Iron', 3),
    (15, 'Bronze', 1), (20, 'Bronze', 2), (25, 'Bronze', 3),
    (30, 'Silver', 1), (35, 'Silver', 2), (40, 'Silver', 3),
    (45, 'Gold', 1), (50, 'Gold', 2), (55, 'Gold', 3),
    (60, 'Platinum', 1), (65, 'Platinum', 2), (70, 'Platinum', 3),
    (75, 'Diamond', 1), (80, 'Diamond', 2), (85, 'Diamond', 3),
    (90, 'Ascendant', 1), (93, 'Ascendant', 2), (96, 'Ascendant', 3),
    (99, 'Immortal', 1), (102, 'Immortal', 2), (105, 'Immortal', 3),
    (110, 'Radiant', 0);
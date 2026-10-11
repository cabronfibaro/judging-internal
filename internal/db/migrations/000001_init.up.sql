CREATE TABLE tracks (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    config TEXT,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE tournaments (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    size INTEGER NOT NULL,
    starting_date INTEGER NOT NULL,
    track_id TEXT NOT NULL REFERENCES tracks(id),
    qualifying_order TEXT NOT NULL,
    randomize_qualifying_order INTEGER NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT 'qualifying'
        CHECK (status IN ('qualifying', 'elimination')),
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    CONSTRAINT size_power_of_2
        CHECK (size > 0 AND (size & (size - 1)) = 0)
);

CREATE TABLE participants (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    challonge_participant_id TEXT,
    tournament_id TEXT NOT NULL
        REFERENCES tournaments(id) ON DELETE CASCADE,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE entrants (
    id TEXT PRIMARY KEY NOT NULL,
    tournament_id TEXT NOT NULL
        REFERENCES tournaments(id) ON DELETE CASCADE,
    participant_id TEXT NOT NULL REFERENCES participants(id),
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('qualified', 'not_qualified', 'pending', 'scored')),
    score REAL NOT NULL DEFAULT 0,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE brackets (
    id TEXT PRIMARY KEY NOT NULL,
    tournament_id TEXT NOT NULL
        REFERENCES tournaments(id) ON DELETE CASCADE,
    challonge_id TEXT NOT NULL,
    url TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'not_started'
        CHECK (status IN ('not_started', 'started')),
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE nodes (
    id TEXT PRIMARY KEY NOT NULL,
    bracket_id TEXT NOT NULL
        REFERENCES brackets(id) ON DELETE CASCADE,
    participant_id TEXT
        REFERENCES participants(id) ON DELETE CASCADE,
    match_id TEXT NOT NULL,
    winner_id TEXT,
    status TEXT NOT NULL,
    round INTEGER NOT NULL,
    score INTEGER DEFAULT 0
);

CREATE TABLE active_battles (
    id TEXT PRIMARY KEY NOT NULL,
    match_id TEXT NOT NULL,
    participant_id TEXT NOT NULL
        REFERENCES participants(id) ON DELETE CASCADE,
    lead_status INTEGER NOT NULL DEFAULT 0,
    rerun INTEGER DEFAULT 0
);

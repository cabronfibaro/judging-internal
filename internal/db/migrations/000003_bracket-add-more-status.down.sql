CREATE TABLE brackets_old (
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

INSERT INTO brackets_old (
    id,
    tournament_id,
    challonge_id,
    url,
    status,
    created_at,
    updated_at
)
SELECT
    id,
    tournament_id,
    challonge_id,
    url,
    status,
    created_at,
    updated_at
FROM brackets;

DROP TABLE brackets;

ALTER TABLE brackets_old RENAME TO brackets;


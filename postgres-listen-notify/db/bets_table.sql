CREATE TABLE IF NOT EXISTS bets (
    timestamp TIMESTAMPTZ NOT NULL,
    username TEXT NOT NULL,
    team TEXT NOT NULL,
    amount FLOAT NOT NULL
);

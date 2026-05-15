CREATE SCHEMA IF NOT EXISTS taskmgr;

CREATE TABLE IF NOT EXISTS taskmgr.users(
    id              SERIAL                      PRIMARY KEY,
    version         BIGINT          NOT NULL    DEFAULT 1,
    name            VARCHAR(100)    NOT NULL                    CHECK (char_length(name) BETWEEN 3 AND 100),
    phone_number    VARCHAR(15)     NOT NULL    UNIQUE          CHECK (char_length(phone_number) BETWEEN 10 AND 15 AND phone_number ~ '^[0-9]+$')
);

CREATE TABLE IF NOT EXISTS taskmgr.tasks(
    id              SERIAL                      PRIMARY KEY,
    version         BIGINT          NOT NULL    DEFAULT 1,
    title           VARCHAR(100)    NOT NULL                    CHECK (char_length(title) BETWEEN 1 AND 100),
    description     VARCHAR(1000)   NOT NULL                    CHECK (char_length(description) BETWEEN 1 AND 1000),
    completed       BOOLEAN         NOT NULL    DEFAULT FALSE,
    created_at      TIMESTAMPTZ     NOT NULL    DEFAULT NOW(),
    completed_at    TIMESTAMPTZ,

    CHECK (
        (completed = FALSE AND completed_at IS NULL) OR (completed = TRUE AND completed_at IS NOT NULL AND completed_at >= created_at)
        ),

    user_id         INT             NOT NULL                    REFERENCES taskmgr.users(id)
);
CREATE SCHEMA todoapp;

CREATE TABLE todoapp.users
(
    id      SERIAL PRIMARY KEY,
    version BIGINT       NOT NULL DEFAULT 1,
    name    VARCHAR(100) NOT NULL CHECK (char_length(name) BETWEEN 3 and 100),
    phone   VARCHAR(15) CHECK (phone ~ '^(?:\+7|8)\d{10}$')
);

CREATE TABLE todoapp.tasks
(
    id           SERIAL PRIMARY KEY,
    version      BIGINT       NOT NULL DEFAULT 1,
    title        VARCHAR(100) NOT NULL,
    description  VARCHAR(1000),
    completed    BOOLEAN               DEFAULT FALSE,
    created_at   TIMESTAMPTZ  NOT NULL,
    completed_at TIMESTAMPTZ,

    CHECK (
        (completed = FALSE AND completed_at IS NULL)
            OR
        (completed = TRUE AND completed_at >= created_at)
        ),

    user_id      INT          NOT NULL REFERENCES todoapp.users (id)
);
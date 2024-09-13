-- +migrate Up
CREATE TABLE IF NOT EXISTS USERS
(
    userId bigint  NOT NULL PRIMARY KEY,
    name   varchar NOT NULL,
    gender varchar NOT NULL,
    age    integer NOT NULL
)
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY not null,
    uuid TEXT not null,
    username TEXT not null,
    password TEXT not null,
    created_at datetime not null,
    updated_at datetime not null,
    confirmed_at datetime not null,
    deleted_at datetime not null
);
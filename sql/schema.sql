CREATE TABLE IF NOT EXISTS passwords (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_name TEXT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    initial_vector TEXT UNIQUE NOT NULL,
    created_at TEXT NULL,
    updated_at TEXT NULL
) STRICT;


CREATE TABLE IF NOT EXISTS master (
   	id INTEGER PRIMARY KEY AUTOINCREMENT,
    password TEXT UNIQUE NOT NULL,
    salt TEXT UNIQUE NOT NULL,
    initial_vector TEXT UNIQUE NOT NULL,
    created_at TEXT NULL,
    updated_at TEXT NULL
) STRICT;
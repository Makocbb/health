CREATE TABLE IF NOT EXISTS ingested_forms (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    session_id TEXT NOT NULL,
    application_reference TEXT NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    gender TEXT NOT NULL,
    date_of_birth TEXT NOT NULL,
    phone_number TEXT,
    mobile_number TEXT NOT NULL,
    address TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS transformed_forms (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    sent_to_bot BOOLEAN NOT NULL DEFAULT FALSE,
    session_id TEXT NOT NULL,
    application_reference TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    gender TEXT NOT NULL,
    date_of_birth TEXT NOT NULL,
    phone_number TEXT,
    mobile_number TEXT NOT NULL,
    address_line_1 TEXT NOT NULL,
    address_line_2 TEXT NOT NULL,
    address_line_3 TEXT,
    postcode TEXT NOT NULL,
    country TEXT NOT NULL,
    longitude REAL NOT NULL,
    latitude REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS transform_logs (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    success INTEGER NOT NULL,
    status TEXT NOT NULL,
    message TEXT,
    transformed_form_id INTEGER NOT NULL
);

CREATE TABLE Characters (
    id SERIAL PRIMARY KEY,
    ID_character INTEGER NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT
);

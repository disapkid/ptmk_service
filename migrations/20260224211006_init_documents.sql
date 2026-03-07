-- +goose Up

CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS documents(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    doc_type TEXT,
    issue_date DATE,
    end_date DATE,
    doc_path TEXT
);

CREATE TABLE IF NOT EXISTS doc_parties(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    doc_id INTEGER REFERENCES documents(id),
    company_name TEXT,
    name TEXT,
    surname TEXT,
    initials TEXT
);

-- +goose Down
DROP TABLE IF EXISTS doc_parties;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS users;

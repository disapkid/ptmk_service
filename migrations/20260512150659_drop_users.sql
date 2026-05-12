-- +goose Up

ALTER TABLE documents
    DROP CONSTRAINT IF EXISTS documents_user_id_fkey;

ALTER TABLE documents
    ALTER COLUMN user_id TYPE BIGINT,
    ALTER COLUMN user_id SET NOT NULL;

DROP TABLE IF EXISTS users;

-- +goose Down

CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY
);

INSERT INTO users (id) OVERRIDING SYSTEM VALUE
SELECT DISTINCT user_id
FROM documents
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('users', 'id'),
    COALESCE((SELECT MAX(id) + 1 FROM users), 1),
    false
);

ALTER TABLE documents
    ADD CONSTRAINT documents_user_id_fkey
    FOREIGN KEY (user_id)
    REFERENCES users(id);

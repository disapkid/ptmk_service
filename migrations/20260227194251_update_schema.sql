-- +goose Up

ALTER TABLE documents
    ALTER COLUMN user_id TYPE BIGINT;

ALTER TABLE doc_parties
    ALTER COLUMN doc_id TYPE BIGINT;

ALTER TABLE documents
    ALTER COLUMN user_id SET NOT NULL;

ALTER TABLE doc_parties
    ALTER COLUMN doc_id SET NOT NULL;

ALTER TABLE documents
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT now();

ALTER TABLE doc_parties
    ADD COLUMN IF NOT EXISTS party_type TEXT;

UPDATE doc_parties
SET party_type = CASE
    WHEN company_name IS NOT NULL THEN 'legal'
    ELSE 'natural'
END
WHERE party_type IS NULL;

ALTER TABLE doc_parties
    ALTER COLUMN party_type SET NOT NULL;

ALTER TABLE doc_parties
    ADD CONSTRAINT doc_parties_party_type_check
    CHECK (party_type IN ('legal','natural'));

ALTER TABLE doc_parties
    RENAME COLUMN name TO first_name;

ALTER TABLE doc_parties
    RENAME COLUMN surname TO last_name;

ALTER TABLE doc_parties
    ADD COLUMN IF NOT EXISTS middle_name TEXT;

CREATE INDEX IF NOT EXISTS idx_documents_user_id
    ON documents(user_id);

CREATE INDEX IF NOT EXISTS idx_doc_parties_doc_id
    ON doc_parties(doc_id);

-- +goose Down

DROP INDEX IF EXISTS idx_doc_parties_doc_id;
DROP INDEX IF EXISTS idx_documents_user_id;

ALTER TABLE doc_parties
    DROP CONSTRAINT IF EXISTS doc_parties_party_type_check;

ALTER TABLE doc_parties
    DROP COLUMN IF EXISTS middle_name;

ALTER TABLE doc_parties
    RENAME COLUMN first_name TO name;

ALTER TABLE doc_parties
    RENAME COLUMN last_name TO surname;

ALTER TABLE doc_parties
    DROP COLUMN IF EXISTS party_type;

ALTER TABLE documents
    DROP COLUMN IF EXISTS created_at;

ALTER TABLE doc_parties
    ALTER COLUMN doc_id TYPE INTEGER;

ALTER TABLE documents
    ALTER COLUMN user_id TYPE INTEGER;
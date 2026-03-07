-- +goose Up

ALTER TABLE doc_parties
    DROP CONSTRAINT IF EXISTS doc_parties_doc_id_fkey;

ALTER TABLE doc_parties
    ADD CONSTRAINT doc_parties_doc_id_fkey
    FOREIGN KEY (doc_id)
    REFERENCES documents(id)
    ON DELETE CASCADE;

-- +goose Down

ALTER TABLE doc_parties
    DROP CONSTRAINT IF EXISTS doc_parties_doc_id_fkey;

ALTER TABLE doc_parties
    ADD CONSTRAINT doc_parties_doc_id_fkey
    FOREIGN KEY (doc_id)
    REFERENCES documents(id)
    ON DELETE NO ACTION;

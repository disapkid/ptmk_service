-- +goose Up
ALTER TABLE documents
    ADD COLUMN doc_number TEXT;

-- +goose Down
ALTER TABLE documents
    DROP COLUMN doc_number;

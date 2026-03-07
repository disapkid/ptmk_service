-- +goose Up
ALTER TABLE documents
    DROP COLUMN created_at;


-- +goose Down
ALTER TABLE documents   
    ADD COLUMN created_at timestamp;

-- +goose up
ALTER TABLE documents RENAME COLUMN end_date TO expire_date;
ALTER TABLE documents RENAME COLUMN issue_date TO start_date;

-- +goose down
ALTER TABLE documents RENAME COLUMN expire_date TO end_date;
ALTER TABLE documents RENAME COLUMN start_date TO issue_date;

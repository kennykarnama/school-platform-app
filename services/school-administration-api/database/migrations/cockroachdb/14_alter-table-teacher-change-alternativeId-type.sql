-- +goose Up
ALTER TABLE teacher ALTER COLUMN alternative_id TYPE STRING;

-- This conversion is intentionally irreversible because existing string IDs
-- are not guaranteed to be valid UUIDs.

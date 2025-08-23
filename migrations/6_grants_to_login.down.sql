ALTER TABLE user_document_access DROP COLUMN IF EXISTS user_login;

ALTER TABLE user_document_access ADD COLUMN user_id INTEGER NOT NULL REFERENCES "user"(id) ON DELETE CASCADE;
ALTER TABLE user_document_access DROP COLUMN IF EXISTS user_id;

ALTER TABLE user_document_access ADD COLUMN user_login VARCHAR(255) NOT NULL REFERENCES "user"(login) ON DELETE CASCADE;
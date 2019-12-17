ALTER TABLE credentials DROP COLUMN IF EXISTS email;
DROP INDEX IF EXISTS credentials_email_idx;

ALTER TABLE credentials DROP COLUMN IF EXISTS phone;
DROP INDEX IF EXISTS credentials_phone_idx;

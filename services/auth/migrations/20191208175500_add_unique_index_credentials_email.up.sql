DROP INDEX IF EXISTS credentials_email_idx;
CREATE UNIQUE INDEX IF NOT EXISTS credentials_email_idx ON credentials(email);

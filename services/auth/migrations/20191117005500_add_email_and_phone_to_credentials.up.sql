ALTER TABLE credentials ADD COLUMN IF NOT EXISTS email varchar(255);
CREATE INDEX IF NOT EXISTS credentials_email_idx on credentials(email);

ALTER TABLE credentials ADD COLUMN IF NOT EXISTS phone varchar(20);
CREATE INDEX IF NOT EXISTS credentials_phone_idx on credentials(phone);

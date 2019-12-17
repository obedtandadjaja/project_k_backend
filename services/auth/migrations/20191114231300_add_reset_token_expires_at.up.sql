ALTER TABLE credentials ADD COLUMN IF NOT EXISTS password_reset_token_expires_at timestamp;

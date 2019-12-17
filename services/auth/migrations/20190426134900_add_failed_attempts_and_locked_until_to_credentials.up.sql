ALTER TABLE credentials ADD COLUMN IF NOT EXISTS failed_attempts integer DEFAULT 0;
ALTER TABLE credentials ADD COLUMN IF NOT EXISTS locked_until timestamp;

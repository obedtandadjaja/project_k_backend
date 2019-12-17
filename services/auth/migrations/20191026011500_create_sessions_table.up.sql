CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SEQUENCE IF NOT EXISTS sessions_id_seq;
CREATE TABLE IF NOT EXISTS sessions(
  id               int PRIMARY KEY DEFAULT nextval('sessions_id_seq'),
  uuid             uuid DEFAULT uuid_generate_v4(),
  credential_id    int REFERENCES credentials(id) ON DELETE CASCADE,
  ip_address       varchar(100),
  user_agent       varchar(255),
  last_accessed_at timestamp DEFAULT now(),
  created_at       timestamp DEFAULT now(),
  expires_at       timestamp DEFAULT now()
);
ALTER SEQUENCE sessions_id_seq OWNED BY sessions.id;
CREATE INDEX IF NOT EXISTS sessions_uuid_idx ON sessions(uuid);

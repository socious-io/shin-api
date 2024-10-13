CREATE TYPE integration_key_status AS ENUM ('ACTIVE', 'SUSPENDED');

CREATE TABLE integration_keys (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  user_id UUID NOT NULL,
  base_url TEXT NOT NULL,
  key TEXT NOT NULL,
  secret TEXT NOT NULL,
  status integration_key_status NOT NULL DEFAULT 'ACTIVE',
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TYPE import_target AS ENUM ('CREDENTIALS');

CREATE TABLE imports (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  user_id UUID NOT NULL,
  target import_target NOT NULL,
  entities UUID[] DEFAULT [],
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
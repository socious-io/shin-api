CREATE TYPE import_target AS ENUM ('CREDENTIALS');
CREATE TYPE import_status AS ENUM ('INITIATED', 'COMPLETED');

CREATE TABLE imports (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  user_id UUID NOT NULL,
  target import_target NOT NULL,
  entities UUID[] DEFAULT '{}',
  count int NOT NULL DEFAULT 0,
  total_count int NOT NULL,
  status import_status NOT NULL DEFAULT 'INITIATED',
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

ALTER TABLE credentials
ADD COLUMN sent BOOLEAN DEFAULT false;
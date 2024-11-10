CREATE TYPE csv_import_doc_type AS ENUM ('CREDENTIALS');
CREATE TYPE csv_import_status AS ENUM ('DONE', 'FAILED');

CREATE TABLE csv_imports (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  user_id UUID NOT NULL,
  doc_type csv_import_doc_type NOT NULL,
  data jsonb DEFAULT NULL,
  status csv_import_status NOT NULL,
  reason text DEFAULT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
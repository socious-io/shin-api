ALTER TABLE recipients ADD COLUMN customer_id VARCHAR(256);
ALTER TABLE recipients ALTER COLUMN email DROP NOT NULL;

CREATE UNIQUE INDEX recipient_user_customer ON recipients (user_id, customer_id) WHERE customer_id IS NOT NULL;

CREATE TABLE credential_verification_individuals (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  recipient_id UUID,
  verification_id UUID,
  user_id UUID,
  status verification_status_type NOT NULL DEFAULT 'CREATED',
  body JSONB,
  validation_error TEXT,
  present_id TEXT,
  connection_id TEXT,
  connection_url TEXT,
  connection_at TIMESTAMP,
  verified_at TIMESTAMP,
  updated_at TIMESTAMP DEFAULT NOW(),
  created_at TIMESTAMP DEFAULT NOW(),  
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_recipient FOREIGN KEY (recipient_id) REFERENCES recipients(id) ON DELETE SET NULL,
  CONSTRAINT fk_verification FOREIGN KEY (verification_id) REFERENCES credential_verifications(id) ON DELETE CASCADE
);

ALTER TABLE credential_verifications 
  DROP COLUMN body,
  DROP COLUMN status,
  DROP COLUMN connection_at,
  DROP COLUMN connection_id,
  DROP COLUMN connection_url,
  DROP COLUMN verified_at,
  DROP COLUMN present_id,
  DROP COLUMN validation_error;
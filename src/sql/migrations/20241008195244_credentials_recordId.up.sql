ALTER TYPE credential_status_type ADD VALUE 'REVOKED';

ALTER TABLE credentials 
  ADD COLUMN record_id UUID,
  ADD COLUMN revoked_at TIMESTAMP;
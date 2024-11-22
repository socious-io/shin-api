CREATE TYPE verification_type AS ENUM ('SINGLE', 'MULTI');

ALTER TABLE credential_verifications ADD COLUMN type verification_type DEFAULT 'MULTI';
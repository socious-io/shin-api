ALTER TABLE credential_schemas ADD COLUMN issue_disabled BOOLEAN DEFAULT false;

UPDATE credential_schemas SET issue_disabled=true WHERE name='KYC';
-- Changing Educational Certificate
UPDATE credential_attributes ca
SET type='DATETIME'
WHERE name='date_of_birth' AND ca.schema_id IN (SELECT id FROM credential_schemas WHERE name='Educational Certificate' AND public=true);

UPDATE credential_attributes ca
SET type='TEXT'
WHERE name='grade' AND ca.schema_id IN (SELECT id FROM credential_schemas WHERE name='Educational Certificate' AND public=true);


-- Changing Work Certificate
UPDATE credential_attributes ca
SET type='DATETIME'
WHERE name='date_of_birth' AND ca.schema_id IN (SELECT id FROM credential_schemas WHERE name='Work Certificate' AND public=true);

DELETE
FROM credential_attributes ca
WHERE name='Company' AND ca.schema_id IN (SELECT id FROM credential_schemas WHERE name='Work Certificate' AND public=true);

-- Updating KYC date_of_birth attribute
UPDATE credential_attributes ca
SET type='DATETIME'
WHERE name='date_of_birth' AND ca.schema_id IN (SELECT id FROM credential_schemas WHERE name='KYC' AND public=true);
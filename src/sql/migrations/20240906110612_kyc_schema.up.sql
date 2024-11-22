-- Insert into credential_schemas and return the id
WITH inserted_schema AS (
  INSERT INTO credential_schemas (name, description, public, deleteable)
  VALUES ('KYC', 'Know your customer schema', true, false)
  RETURNING id
)
-- Insert into credential_attributes using the returned id from the first insert
INSERT INTO credential_attributes (name, description, type, schema_id)
VALUES
  ('first_name', 'First Name', 'TEXT', (SELECT id FROM inserted_schema)),
  ('last_name', 'Last Name', 'TEXT', (SELECT id FROM inserted_schema)),
  ('gender', 'Customer Gender', 'TEXT', (SELECT id FROM inserted_schema)),
  ('id_number', 'provided ID number', 'TEXT', (SELECT id FROM inserted_schema)),
  ('date_of_birth', 'Date of birth', 'TEXT', (SELECT id FROM inserted_schema)),
  ('document_type', 'Document type', 'TEXT', (SELECT id FROM inserted_schema)),
  ('document_number', 'Document Number', 'NUMBER', (SELECT id FROM inserted_schema)),
  ('issued_date', 'KYC issued date', 'DATETIME', (SELECT id FROM inserted_schema));

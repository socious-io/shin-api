-- Adding Educational Certificate
WITH educational_certificate_schema AS (
  INSERT INTO credential_schemas (name, description, public, deleteable)
  VALUES ('Educational Certificate', 'Default schema for academic degrees, diplomas, or certifications', true, false)
  RETURNING id
)

INSERT INTO credential_attributes (name, description, type, schema_id)
VALUES
  ('first_name', 'First Name', 'TEXT', (SELECT id FROM educational_certificate_schema)),
  ('last_name', 'Last Name', 'TEXT', (SELECT id FROM educational_certificate_schema)),
  ('date_of_birth', 'Birth Date', 'TEXT', (SELECT id FROM educational_certificate_schema)),
  ('certificate_name', 'Certificate Name', 'TEXT', (SELECT id FROM educational_certificate_schema)),
  ('field_of_study', 'Field of Study', 'TEXT', (SELECT id FROM educational_certificate_schema)),
  ('start_date', 'Start Date', 'DATETIME', (SELECT id FROM educational_certificate_schema)),
  ('end_date', 'End Date', 'DATETIME', (SELECT id FROM educational_certificate_schema)),
  ('grade', 'Grade', 'NUMBER', (SELECT id FROM educational_certificate_schema)),
  ('description', 'Description', 'TEXT', (SELECT id FROM educational_certificate_schema)),
  ('expiration_date', 'Expiration date', 'DATETIME', (SELECT id FROM educational_certificate_schema));

-- Adding Work Certificate
WITH work_certificate_schema AS (
  INSERT INTO credential_schemas (name, description, public, deleteable)
  VALUES ('Work Certificate', 'Default schema for work history', true, false)
  RETURNING id
)


INSERT INTO credential_attributes (name, description, type, schema_id)
VALUES
  ('first_name', 'First Name', 'TEXT', (SELECT id FROM work_certificate_schema)),
  ('last_name', 'Last Name', 'TEXT', (SELECT id FROM work_certificate_schema)),
  ('date_of_birth', 'Birth Date', 'TEXT', (SELECT id FROM work_certificate_schema)),
  ('job_title', 'Job Title', 'TEXT', (SELECT id FROM work_certificate_schema)),
  ('Company', 'Company', 'TEXT', (SELECT id FROM work_certificate_schema)),
  ('location', 'Location', 'TEXT', (SELECT id FROM work_certificate_schema)),
  ('employment_type', 'Employment Type', 'TEXT', (SELECT id FROM work_certificate_schema)),
  ('start_date', 'Start Date', 'DATETIME', (SELECT id FROM work_certificate_schema)),
  ('end_date', 'End Date', 'DATETIME', (SELECT id FROM work_certificate_schema)),
  ('description', 'Description', 'TEXT', (SELECT id FROM work_certificate_schema));

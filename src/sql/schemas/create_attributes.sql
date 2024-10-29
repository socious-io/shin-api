INSERT INTO credential_attributes (name, description, schema_id, type) 
VALUES (:name, :description, :schema_id, :type)
RETURNING id
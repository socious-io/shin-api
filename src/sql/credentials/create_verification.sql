INSERT INTO credential_verifications (name, description, user_id, schema_id, status) VALUES ($1, $2, $3, $4, 'CREATED')
RETURNING *
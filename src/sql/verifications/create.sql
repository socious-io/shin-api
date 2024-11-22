INSERT INTO credential_verifications (name, description, user_id, schema_id, type) VALUES ($1, $2, $3, $4, $5)
RETURNING *
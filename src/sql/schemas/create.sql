INSERT INTO credential_schemas (name, description, created_id, public) VALUES (
  $1, $2, $3, $4
)
RETURNING *
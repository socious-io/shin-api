INSERT INTO organizations (
  id, name, description, logo
) VALUES ( $1, $2, $3, $4)
RETURNING *
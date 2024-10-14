UPDATE integration_keys
SET name=$2
WHERE id=$1
RETURNING *
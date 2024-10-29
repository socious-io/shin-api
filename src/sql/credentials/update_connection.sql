UPDATE credentials SET
  connection_id=$2,
  connection_url=$3,
  connection_at=NOW(),
  updated_at=NOW()
WHERE id=$1
RETURNING *
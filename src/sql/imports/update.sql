UPDATE imports
SET entities=$2
WHERE id=$1
RETURNING *
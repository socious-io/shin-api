UPDATE recipients SET
  first_name=$2,
  last_name=$3,
  email=$4
WHERE id=$1
RETURNING *
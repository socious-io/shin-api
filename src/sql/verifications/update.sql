update credential_verifications SET
  name=$2,
  description=$3,
  user_id=$4,
  schema_id=$5
WHERE id=$1
RETURNING *
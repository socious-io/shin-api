UPDATE credential_verification_individuals SET
  body=$2,
  validation_error=$3,
  status='FAILED',
  verified_at=NOW(),
  updated_at=NOW()
WHERE id=$1
RETURNING *
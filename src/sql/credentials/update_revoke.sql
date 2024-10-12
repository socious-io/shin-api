UPDATE credentials SET
  status='REVOKED',
  revoked_at=NOW()
WHERE id=$1
RETURNING *
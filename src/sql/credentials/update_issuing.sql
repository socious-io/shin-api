UPDATE credentials SET
  status='ISSUED',
  record_id=$2,
  issued_at=NOW()
WHERE id=$1
RETURNING *
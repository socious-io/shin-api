update organizations SET
  did=$2,
  updated_at=NOW()
WHERE id=$1
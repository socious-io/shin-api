update organizations SET
  name=$2,
  description=$3,
  logo=$4,
  updated_at=NOW()
WHERE id=$1
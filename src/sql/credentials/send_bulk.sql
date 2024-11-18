UPDATE credentials
SET sent=true
WHERE id=ANY($1) AND created_id=$2
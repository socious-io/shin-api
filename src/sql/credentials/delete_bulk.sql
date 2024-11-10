DELETE
FROM credentials
WHERE id=ANY($1) AND created_id=$2
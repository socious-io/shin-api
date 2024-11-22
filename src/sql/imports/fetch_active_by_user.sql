SELECT *
FROM imports
WHERE user_id=$1 AND status='INITIATED'
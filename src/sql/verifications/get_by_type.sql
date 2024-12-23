SELECT 
  id, COUNT(*) OVER () as total_count
FROM credential_verifications cv
WHERE cv.user_id = $1 AND cv.type=$4 LIMIT $2 OFFSET $3
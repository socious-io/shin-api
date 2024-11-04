SELECT 
  id, COUNT(*) OVER () as total_count
FROM credentials cv
WHERE cv.created_id = $1 AND cv.schema_id=$4 LIMIT $2 OFFSET $3
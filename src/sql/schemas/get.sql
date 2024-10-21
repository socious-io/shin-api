SELECT id, COUNT(*) OVER () as total_count 
FROM credential_schemas 
WHERE (created_id = $1 or public = true) LIMIT $2 OFFSET $3
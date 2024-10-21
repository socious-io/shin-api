SELECT 
    id, 
    COUNT(*) OVER () as total_count 
FROM 
    recipients 
WHERE 
    (first_name ILIKE '%' || $1 || '%' OR last_name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%') 
    AND user_id = $2
LIMIT $3 OFFSET $4;
UPDATE imports
SET 
  entities=array_append(entities, $2),
  count=count+1,
  status=(CASE WHEN total_count=count+1 THEN 'COMPLETED'::import_status ELSE 'INITIATED'::import_status END)
WHERE id = $1
RETURNING *;
SELECT id, COUNT(*) OVER () as total_count
	FROM integration_keys ik
	WHERE ik.user_id=$1
LIMIT $2 OFFSET $3
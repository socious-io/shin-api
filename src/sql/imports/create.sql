INSERT INTO imports(user_id, target, total_count)
VALUES($1, $2, $3)
RETURNING *
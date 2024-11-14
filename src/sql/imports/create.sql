INSERT INTO imports(user_id, target)
VALUES($1, $2)
RETURNING *
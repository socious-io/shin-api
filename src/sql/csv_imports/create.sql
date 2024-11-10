INSERT INTO csv_imports(user_id, doc_type, data, status, reason)
VALUES($1, $2, $3, $4, $5)
RETURNING *
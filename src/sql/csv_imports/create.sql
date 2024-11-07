INSERT INTO csv_imports(user_id, doc_type)
VALUES($1, $2)
RETURNING *
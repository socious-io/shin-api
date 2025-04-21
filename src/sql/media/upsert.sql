INSERT INTO media(id, user_id, url, filename)
VALUES($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE SET
    filename = EXCLUDED.filename,
    url = EXCLUDED.url,
    user_id = EXCLUDED.user_id
RETURNING *;
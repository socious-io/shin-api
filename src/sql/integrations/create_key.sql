INSERT
INTO integration_keys (name, user_id, base_url, key, secret)
VALUES ($1, $2, $3, $4, $5)
RETURNING *
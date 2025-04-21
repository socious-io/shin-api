INSERT INTO oauth_connects (user_id, provider, matrix_unique_id, access_token, refresh_token, meta, expired_at)
VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (matrix_unique_id, provider) 
DO UPDATE SET
    access_token = EXCLUDED.access_token,
    refresh_token = EXCLUDED.refresh_token,
    meta = EXCLUDED.meta,
    expired_at = EXCLUDED.expired_at,
    updated_at = NOW();
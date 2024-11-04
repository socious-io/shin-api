INSERT INTO users (first_name, last_name, username, email, password, password_expired)
VALUES ($1, $2, $3, $4, $5, $5::text IS NULL)
ON CONFLICT (email) 
DO UPDATE SET
  updated_at = NOW()
RETURNING *;
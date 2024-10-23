INSERT INTO recipients (first_name, last_name, email, user_id, customer_id)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, customer_id) 
WHERE customer_id IS NOT NULL
DO UPDATE SET 
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    email = EXCLUDED.email
RETURNING *;
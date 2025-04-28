INSERT INTO users (id, first_name, last_name, username, email, bio, phone, avatar) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id) DO UPDATE SET
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    username = EXCLUDED.username,
    email = EXCLUDED.email,
    bio = EXCLUDED.bio,
    phone = EXCLUDED.phone,
    avatar = EXCLUDED.avatar
RETURNING *;
INSERT INTO organizations (
  id, name, description, logo, verified
) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id)
DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    logo=EXCLUDED.logo,
    verified=EXCLUDED.verified
RETURNING *;
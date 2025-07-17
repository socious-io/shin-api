INSERT INTO organization_members (
  user_id, organization_id
) VALUES ($1, $2)
ON CONFLICT (user_id, organization_id) DO NOTHING
RETURNING *
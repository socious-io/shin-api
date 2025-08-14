SELECT om.organization_id AS id
FROM organization_members om
WHERE user_id=$1
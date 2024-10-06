SELECT o.*,
m.url as "logo.url",
m.filename "logo.filename",
(SELECT status FROM kyb_verifications kv WHERE o.id = kv.organization_id ORDER BY created_at DESC LIMIT 1) AS verification_status
FROM organizations o
JOIN organization_members om ON user_id=$1 AND om.organization_id=o.id
LEFT JOIN media m ON o.logo_id=m.id
ORDER BY o.created_at ASC
LIMIT 1
SELECT o.*,
(SELECT status FROM kyb_verifications kv WHERE o.id = kv.organization_id ORDER BY created_at DESC LIMIT 1) AS verification_status
FROM organizations o
JOIN organization_members om ON user_id=$2 AND om.organization_id=o.id
WHERE o.id=$1
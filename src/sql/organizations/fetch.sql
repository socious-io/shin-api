SELECT o.*,
(SELECT status FROM kyb_verifications kv WHERE o.id = kv.organization_id ORDER BY created_at DESC LIMIT 1) AS verification_status
FROM organizations o
WHERE o.id IN(?)
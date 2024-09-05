SELECT o.*,
m.url as "logo.url",
m.filename "logo.filename"
FROM organizations o
JOIN organization_members om ON user_id=$2 AND om.organization_id=o.id
LEFT JOIN media m ON o.logo_id=m.id
WHERE o.id=$1
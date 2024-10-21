SELECT 
  cv.*,
  row_to_json(u.*) AS user,
  row_to_json(r.*) AS recipient,
  row_to_json(v.*) AS verification
FROM credential_verification_individuals cv 
LEFT JOIN users u ON u.id = cv.user_id
LEFT JOIN recipients r ON r.id = cv.recipient_id
LEFT JOIN credential_verifications v ON v.id = cv.verification_id
WHERE cv.id IN (?)
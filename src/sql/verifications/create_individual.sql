INSERT INTO credential_verification_individuals (user_id, recipient_id, verification_id) VALUES ($1, $2, $3)
ON CONFLICT (recipient_id, verification_id) DO UPDATE SET updated_at = NOW()
RETURNING *


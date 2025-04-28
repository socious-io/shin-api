ALTER TABLE users ADD COLUMN avatar jsonb;
ALTER TABLE organizations ADD COLUMN logo jsonb;

UPDATE users
SET avatar = row_to_json(m.*)
FROM media m
WHERE users.avatar_id = m.id;

UPDATE organizations
SET logo = row_to_json(m.*)
FROM media m
WHERE organizations.logo_id = m.id;

ALTER TABLE users DROP COLUMN avatar_id;
ALTER TABLE organizations DROP COLUMN logo_id;

CREATE TABLE oauth_connects (
    id UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    provider TEXT NOT NULL,
    matrix_unique_id TEXT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    meta JSONB,
    expired_at TIMESTAMP,
    created_at TIMESTAMP  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX unique_provider_mui
ON oauth_connects (matrix_unique_id, provider);

ALTER TABLE organizations
RENAME COLUMN is_verified TO verified;
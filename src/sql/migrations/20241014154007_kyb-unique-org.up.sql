-- Delete the least recent KYBs and remain one last recent one per organization (newest remains)
DELETE FROM kyb_verifications
WHERE (organization_id, created_at) NOT IN (
    SELECT organization_id, MAX(created_at)
    FROM kyb_verifications
    GROUP BY organization_id
);

-- Making organization_id unique among KYB columns
ALTER TABLE kyb_verifications ADD CONSTRAINT kyb_verifications_unique_organization_id UNIQUE (organization_id);
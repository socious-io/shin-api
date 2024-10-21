SELECT 
  cs.*,
  row_to_json(u.*) AS created,
  (SELECT
      jsonb_agg(json_build_object(
          'id', id,
          'name', name,
          'description', description,
          'type', type,
          'created_at', created_at
        ))
        FROM credential_attributes ca
        WHERE ca.schema_id=cs.id
    ) AS attributes
FROM credential_schemas cs 
LEFT JOIN users u ON u.id=cs.created_id
WHERE cs.id IN(?)
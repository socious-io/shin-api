SELECT 
  cv.*,
  row_to_json(u.*) AS user,
  row_to_json(cs.*) AS schema,
  row_to_json(vi.*) AS single,
  (SELECT
      jsonb_agg(json_build_object(
          'id', id,
          'attribute_id', attribute_id,
          'schema_id', schema_id,
          'verification_id', verification_id,
          'value', value,
          'operator', operator,
          'created_at', created_at
        ))
        FROM verification_attribute_values va
        WHERE va.verification_id=cv.id
    ) AS attributes
FROM credential_verifications cv 
LEFT JOIN users u ON u.id = cv.user_id
LEFT JOIN credential_schemas cs ON cs.id = cv.schema_id
LEFT JOIN credential_verification_individuals vi ON vi.verification_id = cv.id AND cv.type = 'SINGLE'
WHERE cv.id IN (?)
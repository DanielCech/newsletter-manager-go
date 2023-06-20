UPDATE
    "newsletter"
SET
    name = @name,
    description = @description,
    created_at = @created_at,
    updated_at = @updated_at
WHERE
    id = @id

UPDATE
    "author"
SET
    name = @name,
    email = @email,
    password_hash = @password_hash,
    role = @role,
    created_at = @created_at,
    updated_at = @updated_at
WHERE
    id = @id

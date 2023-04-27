UPDATE
    "user"
SET
    referrer_id = @referrer_id,
    name = @name,
    email = @email,
    password_hash = @password_hash,
    role = @role,
    created_at = @created_at,
    updated_at = @updated_at
WHERE
    id = @id

INSERT INTO "user" (
    id,
    referrer_id,
    name,
    email,
    password_hash,
    role,
    created_at,
    updated_at
)
VALUES
    (@id, @referrer_id, @name, @email, @password_hash, @role, @created_at, @updated_at)

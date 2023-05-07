INSERT INTO "author" (
    id,
    name,
    email,
    password_hash,
    created_at,
    updated_at
)
VALUES
    (@id, @name, @email, @password_hash, @created_at, @updated_at)

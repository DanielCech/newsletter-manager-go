SELECT
    u.id,
    u.name,
    u.email,
    u.password_hash,
    u.created_at,
    u.updated_at
FROM
    "author" AS u

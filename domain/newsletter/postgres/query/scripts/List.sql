SELECT
    u.id,
    u.referrer_id,
    u.name,
    u.email,
    u.password_hash,
    u.role,
    u.created_at,
    u.updated_at
FROM
    "user" AS u

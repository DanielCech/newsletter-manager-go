SELECT
    rt.id,
    rt.user_id,
    rt.expires_at,
    rt.created_at
FROM
    refresh_token AS rt
WHERE
    rt.id = @id
FOR UPDATE

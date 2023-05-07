SELECT
    rt.id,
    rt.author_id,
    rt.expires_at,
    rt.created_at
FROM
    refresh_token AS rt
WHERE
    rt.id = @id
FOR UPDATE

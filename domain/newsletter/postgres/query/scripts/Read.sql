SELECT
    n.id,
    n.author_id,
    n.name,
    n.description,
    n.created_at,
    n.updated_at
FROM
    "newsletter" AS n
WHERE
    n.id = @id
FOR UPDATE
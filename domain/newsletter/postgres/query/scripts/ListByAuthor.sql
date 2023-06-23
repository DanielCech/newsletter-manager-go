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
    n.author_id = @author_id
FOR UPDATE
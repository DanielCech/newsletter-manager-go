INSERT INTO "newsletter" (
    id,
    author_id,
    name,
    description,
    created_at,
    updated_at
)
VALUES
    (@id, @author_id, @name, @description, @created_at, @updated_at)

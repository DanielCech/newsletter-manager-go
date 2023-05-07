INSERT INTO refresh_token (
    id,
    author_id,
    expires_at,
    created_at
)
VALUES
    (@id, @author_id, @expires_at, @created_at)

INSERT INTO refresh_token (
    id,
    user_id,
    expires_at,
    created_at
)
VALUES
    (@id, @user_id, @expires_at, @created_at)

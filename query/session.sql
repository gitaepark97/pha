-- name: CreateSession :exec
INSERT INTO session(
  id,
  user_id,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expired_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
);

-- name: GetSession :one
SELECT
  *
FROM session
WHERE id = ?;
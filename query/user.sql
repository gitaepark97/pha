-- name: CreateUser :exec
INSERT INTO user(
  phone_number,
  hashed_password
) VALUES (
  ?, ?
);

-- name: GetUser :one
SELECT
  *
FROM user
WHERE phone_number = ?;
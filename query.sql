-- name: GetUsername :one
SELECT username FROM users
WHERE username = $1 LIMIT 1;

-- name: GetPassword :one
SELECT password FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserById :one
SELECT name, username, email FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, name, username FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (name, username, email, password)
VALUES ($1, $2, $3, $4);

-- name: UpdateUser :exec
UPDATE users SET
    name = $1,
    username = $2,
    email = $3,
    password = $4
WHERE id = $5;


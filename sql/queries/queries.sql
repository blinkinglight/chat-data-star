-- name: UserByEmail :one
SELECT
    *
FROM
    users
WHERE
    username = @email;

-- name: UserById :one
SELECT
    *
FROM
    users
WHERE
    id = @id;
-- name: GetURLByID :one
SELECT * FROM urls WHERE id = $1;

-- name: ListURLs :many
SELECT * FROM urls ORDER BY created_at DESC LIMIT $1 OFFSET $2;
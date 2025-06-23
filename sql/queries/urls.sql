-- name: GetURLByID :one
SELECT * FROM urls WHERE id = $1;

-- name: ListURLs :many
SELECT * FROM urls ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CreateURL :one
INSERT INTO urls (
    url, frequency, status, max_retries, timeout, rate_limit, 
    user_agent, parser_config, next_scrape_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;
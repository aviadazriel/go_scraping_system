-- name: GetURLByID :one
SELECT * FROM urls WHERE id = $1;

-- name: ListURLs :many
SELECT * FROM urls ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountURLs :one
SELECT COUNT(*) FROM urls;

-- name: CreateURL :one
INSERT INTO urls (
    url, frequency, status, max_retries, timeout, rate_limit, 
    user_agent, parser_config, next_scrape_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetURLsScheduledForScraping :many
SELECT * FROM urls 
WHERE next_scrape_at BETWEEN $1 AND $2 
AND status IN ('pending', 'retry')
ORDER BY next_scrape_at ASC 
LIMIT $3;

-- name: GetURLsByStatus :many
SELECT * FROM urls 
WHERE status = $1 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3;

-- name: UpdateURLStatus :exec
UPDATE urls SET status = $2, updated_at = NOW() WHERE id = $1;

-- name: UpdateNextScrapeTime :exec
UPDATE urls SET next_scrape_at = $2, updated_at = NOW() WHERE id = $1;

-- name: UpdateLastScrapedTime :exec
UPDATE urls SET last_scraped_at = $2, updated_at = NOW() WHERE id = $1;

-- name: IncrementRetryCount :exec
UPDATE urls SET retry_count = retry_count + 1, updated_at = NOW() WHERE id = $1;

-- name: ResetRetryCount :exec
UPDATE urls SET retry_count = 0, updated_at = NOW() WHERE id = $1;

-- name: GetURLsForImmediateScraping :many
SELECT * FROM urls 
WHERE next_scrape_at <= $1 
AND status IN ('pending', 'retry')
ORDER BY next_scrape_at ASC 
LIMIT $2;

-- name: CountURLsByStatus :one
SELECT COUNT(*) FROM urls WHERE status = $1;

-- name: GetURLsByIDs :many
SELECT * FROM urls WHERE id = ANY($1::uuid[]);
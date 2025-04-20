-- name: GetPosts :many
SELECT posts.*
FROM posts
LEFT JOIN feeds ON feed_id = feeds.id
WHERE feeds.user_id = $1
ORDER BY published_at DESC
LIMIT $2;
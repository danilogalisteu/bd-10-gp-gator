-- name: GetPosts :many
SELECT posts.*
FROM posts
INNER JOIN feed_follows ON feed_follows.feed_id = posts.feed_id AND feed_follows.user_id = $1
ORDER BY published_at DESC
LIMIT $2;
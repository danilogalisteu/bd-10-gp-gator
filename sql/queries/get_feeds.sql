-- name: GetFeeds :many
SELECT feeds.*, users.name AS user_name
FROM feeds
LEFT JOIN users ON user_id = users.id;
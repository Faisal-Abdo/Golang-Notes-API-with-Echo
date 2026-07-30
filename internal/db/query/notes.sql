-- name: GetAllNotes :many
SELECT id, title, content, created_at
FROM notes
ORDER BY id;
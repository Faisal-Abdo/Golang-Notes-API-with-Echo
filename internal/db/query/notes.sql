-- name: GetAllNotes :many
SELECT id, title, content, created_at
FROM notes
ORDER BY id;

-- name: GetNoteByID :one
SELECT id, title, content, created_at
FROM notes
WHERE id = $1;

-- name: CreateNote :one
INSERT INTO notes (title, content)
VALUES ($1, $2)
RETURNING id, title, content, created_at;

-- name: UpdateNoteByID :one
UPDATE notes
SET title = $1, content = $2
WHERE id = $3
RETURNING id, title, content, created_at;

-- name: DeleteNoteByID :execrows
DELETE FROM notes WHERE id = $1;

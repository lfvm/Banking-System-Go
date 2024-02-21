-- name: CreateEntry :one
INSERT INTO entries (
  account_id, amount
) VALUES (
  $1, $2
)
RETURNING *;


-- name: ListEntries :many
SELECT * FROM transfers
ORDER BY id 
LIMIT $1 OFFSET $2;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;

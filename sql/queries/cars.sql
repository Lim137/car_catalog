-- name: DeleteCarById :exec
DELETE FROM cars
WHERE id = $1;
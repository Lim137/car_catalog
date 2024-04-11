-- name: CreateCar :one
INSERT INTO cars ( id, created_at, updated_at, reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic )
VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: DeleteCarById :exec
DELETE FROM cars
WHERE id = $1;
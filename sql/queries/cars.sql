-- name: CreateCar :one
INSERT INTO cars ( id, created_at, updated_at, reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic )
VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id;

-- name: DeleteCarById :exec
DELETE FROM cars
WHERE id = $1;

-- name: UpdateRegNumById :exec
UPDATE cars
SET reg_num = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateMarkById :exec
UPDATE cars
SET mark = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateModelById :exec
UPDATE cars
SET model = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateYearById :exec
UPDATE cars
SET year = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateOwnerNameById :exec
UPDATE cars
SET owner_name = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateOwnerSurnameById :exec
UPDATE cars
SET owner_surname = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateOwnerPatronymicById :exec
UPDATE cars
SET owner_patronymic = $2, updated_at = NOW()
WHERE id = $1;

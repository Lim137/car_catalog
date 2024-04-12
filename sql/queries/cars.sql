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

-- name: GetCars :many
SELECT *
FROM cars
WHERE
    (reg_num = $1 OR $1 = '') AND
    (mark = $2 OR $2 = '') AND
    (model = $3 OR $3 = '') AND
    (year = $4 OR $4 = 0) AND
    (owner_name = $5 OR $5 = '') AND
    (owner_surname = $6 OR $6 = '') AND
    (owner_patronymic = $7 OR $7 IS NULL );

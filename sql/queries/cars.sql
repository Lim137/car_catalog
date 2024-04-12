-- name: CreateCar :one
INSERT INTO cars ( id, created_at, updated_at, reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic )
VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id;

-- name: DeleteCarById :exec
DELETE FROM cars
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
    (owner_patronymic = $7 OR $7 = '')
LIMIT CASE WHEN $8 = -1 THEN NULL ELSE $8 END
OFFSET ($9 - 1) * $8;

-- name: UpdateCarById :one
UPDATE cars
SET
    updated_at = NOW(),
    reg_num = CASE
                   WHEN $2 != '' THEN $2
                   ELSE reg_num
            END,
    year = CASE
               WHEN $5 != -1 THEN $5
               ELSE year
        END,
    mark = CASE
               WHEN $3 != '' THEN $3
               ELSE mark
        END,
    model = CASE
                WHEN $4 != '' THEN $4
                ELSE model
        END,
    owner_name = CASE
                WHEN $6 != '' THEN $6
                ELSE owner_name
        END,
    owner_surname = CASE
                WHEN $7 != '' THEN $7
                ELSE owner_surname
        END,
    owner_patronymic = CASE
                WHEN $8 != '' THEN $8
                ELSE owner_patronymic
        END
WHERE id = $1
RETURNING *;
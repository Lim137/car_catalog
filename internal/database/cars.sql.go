// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: cars.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createCar = `-- name: CreateCar :one
INSERT INTO cars ( id, created_at, updated_at, reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic )
VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id
`

type CreateCarParams struct {
	ID              uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	RegNum          string
	Mark            string
	Model           string
	Year            int32
	OwnerName       string
	OwnerSurname    string
	OwnerPatronymic string
}

func (q *Queries) CreateCar(ctx context.Context, arg CreateCarParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createCar,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.RegNum,
		arg.Mark,
		arg.Model,
		arg.Year,
		arg.OwnerName,
		arg.OwnerSurname,
		arg.OwnerPatronymic,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteCarById = `-- name: DeleteCarById :exec
DELETE FROM cars
WHERE id = $1
`

func (q *Queries) DeleteCarById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteCarById, id)
	return err
}

const getCars = `-- name: GetCars :many
SELECT id, created_at, updated_at, reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic
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
OFFSET ($9 - 1) * $8
`

type GetCarsParams struct {
	RegNum          string
	Mark            string
	Model           string
	Year            int32
	OwnerName       string
	OwnerSurname    string
	OwnerPatronymic string
	Column8         interface{}
	Column9         interface{}
}

func (q *Queries) GetCars(ctx context.Context, arg GetCarsParams) ([]Car, error) {
	rows, err := q.db.QueryContext(ctx, getCars,
		arg.RegNum,
		arg.Mark,
		arg.Model,
		arg.Year,
		arg.OwnerName,
		arg.OwnerSurname,
		arg.OwnerPatronymic,
		arg.Column8,
		arg.Column9,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Car
	for rows.Next() {
		var i Car
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.RegNum,
			&i.Mark,
			&i.Model,
			&i.Year,
			&i.OwnerName,
			&i.OwnerSurname,
			&i.OwnerPatronymic,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCarById = `-- name: UpdateCarById :one
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
RETURNING id, created_at, updated_at, reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic
`

type UpdateCarByIdParams struct {
	ID      uuid.UUID
	Column2 interface{}
	Column3 interface{}
	Column4 interface{}
	Column5 interface{}
	Column6 interface{}
	Column7 interface{}
	Column8 interface{}
}

func (q *Queries) UpdateCarById(ctx context.Context, arg UpdateCarByIdParams) (Car, error) {
	row := q.db.QueryRowContext(ctx, updateCarById,
		arg.ID,
		arg.Column2,
		arg.Column3,
		arg.Column4,
		arg.Column5,
		arg.Column6,
		arg.Column7,
		arg.Column8,
	)
	var i Car
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RegNum,
		&i.Mark,
		&i.Model,
		&i.Year,
		&i.OwnerName,
		&i.OwnerSurname,
		&i.OwnerPatronymic,
	)
	return i, err
}

-- +goose Up
CREATE TABLE cars (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    reg_num TEXT NOT NULL,
    mark TEXT NOT NULL,
    model TEXT NOT NULL,
    year INT NOT NULL,
    owner_name TEXT NOT NULL,
    owner_surname TEXT NOT NULL,
    owner_patronymic TEXT NOT NULL
);
-- +goose Down
DROP TABLE cars;
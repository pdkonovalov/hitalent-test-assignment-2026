-- +goose Up
CREATE TABLE employees (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    department_id bigint REFERENCES departments(id) ON DELETE CASCADE,
    full_name varchar(200) NOT NULL,
    position varchar(200) NOT NULL,
    hired_at timestamptz,
    created_at timestamptz DEFAULT now()
);

-- +goose Down
DROP TABLE employees;

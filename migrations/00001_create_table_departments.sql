-- +goose Up
CREATE TABLE departments (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name varchar(200) NOT NULL,
    parent_id bigint REFERENCES departments(id) ON DELETE CASCADE,
    created_at timestamptz DEFAULT now(),
    UNIQUE (name, parent_id)
);

-- +goose Down
DROP TABLE departments;

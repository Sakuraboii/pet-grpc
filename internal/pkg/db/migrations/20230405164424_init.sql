-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    name       text                  NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE       DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE TABLE cars
(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    model varchar,
    user_id BIGSERIAL REFERENCES users (id),
    created_at TIMESTAMP WITH TIME ZONE       DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE cars;
-- +goose StatementEnd

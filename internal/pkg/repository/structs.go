package repository

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64        `db:"id"`
	Name      string       `db:"name"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type Car struct {
	ID        int64        `db:"id"`
	Model     string       `db:"model"`
	UserId    int64        `db:"user_id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type Pagination struct {
	Page  int
	Limit int
}

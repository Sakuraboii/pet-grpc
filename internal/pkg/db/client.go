package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
)

const (
	host     = "postgres"
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func NewDB(ctx context.Context) (*sqlx.DB, *Database, error) {
	config, err := pgx.ParseURI(GenerateDsn())
	if err != nil {
		return nil, nil, err
	}
	nativeDB := stdlib.OpenDB(config)
	db := sqlx.NewDb(nativeDB, "pgxConn")

	pgxConn, _ := pgxpool.Connect(ctx, GenerateDsn())

	return db, newDatabase(pgxConn), db.Ping()
}

func GenerateDsn() string {
	return fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname)
}

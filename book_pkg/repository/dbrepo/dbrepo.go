package dbrepo

import (
	"bookings/book_pkg/config"
	"bookings/book_pkg/repository"
	"database/sql"
)

type PostgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}
type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresDBRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &PostgresDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}

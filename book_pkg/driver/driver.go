package driver

import (
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

type DB struct {
	SQL *sql.DB
}

var connDB = &DB{}

func ConnectDBS(dns string) (*DB, error) {
	conn, err := execDB(dns)
	if err != nil {
		panic(err)
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	connDB.SQL = conn

	err = TestDBConn(conn)
	if err != nil {
		return nil, err
	}

	return connDB, nil
}

func TestDBConn(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func execDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return db, nil
	}

	err = db.Ping()
	if err != nil {
		return db, nil
	}

	return db, nil
}

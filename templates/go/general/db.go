package general

import "github.com/victorolegovich/sgen/settings"

const (
	mysql string = `package db

import (
	"github.com/jmoiron/sqlx"
)

const connection = ""

func NewConnection() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", connection)
	if err != nil {
		return nil, err
	}
	return db, nil
}`
	postgres string = `package db

import (
	"context"
	"github.com/jackc/pgx"
)

//Fill in the configuration
const connection string = ""

func NewConnection() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		return nil, err
	}

	return conn, nil
}`
)

func Src(driver string) string {
	switch driver {
	case settings.MySQL:
		return mysql
	case settings.PostgreSQL:
		return postgres
	}
	return ""
}

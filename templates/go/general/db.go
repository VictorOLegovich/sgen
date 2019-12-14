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
	"github.com/jackc/pgx"
)

//Fill in the configuration
func config() (c pgx.ConnConfig) {
	return c
}

func NewConnection() (*pgx.Conn, error) {
	conn, err := pgx.Connect(config())
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

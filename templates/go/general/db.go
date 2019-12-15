package general

import "github.com/victorolegovich/sgen/settings"

const (
	mysql string = `
package db

import (
	"github.com/jmoiron/sqlx"
	"os"
	"strings"
)

type connection struct {
	user, pass, dbname, sslmode string
}

func newConnection(user string, pass string, dbname string, sslmode string) *connection {
	return &connection{user: user, pass: pass, dbname: dbname, sslmode: sslmode}
}

func (c *connection) string() string {
	var builder strings.Builder

	if c.user != "" {
		builder.WriteString("user=")
		builder.WriteString(c.user)
		builder.WriteString(" ")
	} else {
		println(` + "`No \"user\" value is transmitted when connecting to the database `" + `)
		os.Exit(1)
	}

	if c.pass != "" {
		builder.WriteString("password=")
		builder.WriteString(c.pass)
		builder.WriteString(" ")
	}

	if c.dbname != "" {
		builder.WriteString("dbname=")
		builder.WriteString(c.dbname)
		builder.WriteString(" ")
	} else {
		println(` + "`the dbname parameter is not specified, so you should specify it in query_builder if you use it`" + `)
	}

	if c.sslmode != "" {
		builder.WriteString("sslmode=")
		builder.WriteString(c.sslmode)
	}

	return builder.String()
}


func NewConnection() (*sqlx.DB, error) {
	connectionInfo := newConnection("postgres", "", "postgres", "disable").string()
	db, err := sqlx.Connect("mysql", connectionInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}`
	postgres string = `package db

import (
	"context"
	"github.com/jackc/pgx"
	"os"
	"strings"
)


type connection struct {
	user, pass, dbname, sslmode string
}

func newConnection(user string, pass string, dbname string, sslmode string) *connection {
	return &connection{user: user, pass: pass, dbname: dbname, sslmode: sslmode}
}

func (c *connection) string() string {
	var builder strings.Builder

	if c.user != "" {
		builder.WriteString("user=")
		builder.WriteString(c.user)
		builder.WriteString(" ")
	} else {
		println(` + "`No \"user\" value is transmitted when connecting to the database `" + `)
		os.Exit(1)
	}

	if c.pass != "" {
		builder.WriteString("password=")
		builder.WriteString(c.pass)
		builder.WriteString(" ")
	}

	if c.dbname != "" {
		builder.WriteString("dbname=")
		builder.WriteString(c.dbname)
		builder.WriteString(" ")
	} else {
		println(` + "`the dbname parameter is not specified, so you should specify it in query_builder if you use it`" + `)
	}

	if c.sslmode != "" {
		builder.WriteString("sslmode=")
		builder.WriteString(c.sslmode)
	}

	return builder.String()
}

func NewConnection() (*pgx.Conn, error) {
	connectionInfo := newConnection("postgres", "", "postgres", "disable").string()
	conn, err := pgx.Connect(context.Background(), connectionInfo)
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

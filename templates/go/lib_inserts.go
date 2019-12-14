package _go

import (
	"github.com/victorolegovich/sgen/settings"
)

type libInsert struct {
	driver string
}

func newLibInsert(driver string) *libInsert {
	return &libInsert{driver}
}

func (l *libInsert) toImport() string {
	switch l.driver {
	case settings.MySQL:
		return "\t\"github.com/jmoiron/sqlx\"\n" +
			"\t_ \"github.com/go-sql-driver/mysql\"\n"
	case settings.PostgreSQL:
		return "\t\"context\"\n\t\"github.com/jackc/pgx\"\n"
	default:
		return ""
	}
}

func (l *libInsert) toType() string {
	switch l.driver {
	case settings.MySQL:
		return "*sqlx.DB"
	case settings.PostgreSQL:
		return "*pgx.Conn"
	default:
		return ""
	}
}

func (l *libInsert) toReadOne(value, scanning string, lineBreak bool) string {
	pgxROne := "err = s.db.QueryRow(\n\t\tquery," + value + ").Scan("
	if lineBreak {
		return pgxROne + "\n" + scanning + "\t)"
	} else {
		return pgxROne + scanning + ")"
	}
}

func (l *libInsert) toReadAll(fillingVar, addVar, scanning string, lineBreak bool) string {
	pgxRow := "rows, err := s.db.Query(query)\n\n" +
		"\tif err != nil{\n" +
		"\t\treturn nil, err\n" +
		"\t}\n\n" +
		"\tfor rows.Next(){\n" +
		"\t\tif err := rows.Scan("
	if lineBreak {
		pgxRow += "\n"
	}
	pgxRow += scanning + "\t\t); err != nil{\n\t\t\treturn nil, err\n\t\t}\n\n" +
		"\t\t" + addVar + " = append(" + addVar + "," + fillingVar + ")\n\t}"

	return pgxRow
}

func (l *libInsert) toExec(args string, lineBreak bool) string {
	exec := "\t_, err = s.db.Exec("
	if lineBreak {
		exec += "\n\t\tquery," + args + "\t)\n"
	} else {
		exec += "query, " + args + ")\n"
	}
	return exec
}

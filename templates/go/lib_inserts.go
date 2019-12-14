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
	pgx := "err = s.db.QueryRow(context.Background(), query," + value + ").Scan("
	sqlx := "err = s.db.QueryRow( query," + value + ").Scan("
	if lineBreak {
		pgx += "\n" + scanning + "\t)"
		sqlx += "\n" + scanning + "\t)"
	} else {
		pgx += scanning + ")"
		sqlx += scanning + ")"
	}

	switch l.driver {
	case settings.MySQL:
		return sqlx
	case settings.PostgreSQL:
		return pgx
	}
	return ""
}

func (l *libInsert) toReadAll(fillingVar, addVar, scanning string, lineBreak bool) string {
	pgxRow := "rows, err := s.db.Query(context.Background(), query)\n\n" +
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

	sqlxRow := "rows, err := s.db.Query(context.Background(), query)\n\n" +
		"\tif err != nil{\n" +
		"\t\treturn nil, err\n" +
		"\t}\n\n" +
		"\tfor rows.Next(){\n" +
		"\t\tif err := rows.Scan("
	if lineBreak {
		sqlxRow += "\n"
	}
	sqlxRow += scanning + "\t\t); err != nil{\n\t\t\treturn nil, err\n\t\t}\n\n" +
		"\t\t" + addVar + " = append(" + addVar + "," + fillingVar + ")\n\t}"

	switch l.driver {
	case settings.PostgreSQL:
		return pgxRow
	case settings.MySQL:
		return sqlxRow
	}

	return ""
}

func (l *libInsert) toExec(args string, lineBreak bool) string {
	pgx := "\t_, err = s.db.Exec(context.Background(),"
	sqlx := "\t_, err = s.db.Exec("
	if lineBreak {
		pgx += "\n\t\tquery," + args + "\t)\n"
		sqlx += "\n\t\tquery," + args + "\t)\n"
	} else {
		pgx += "query, " + args + ")\n"
		sqlx += "query, " + args + ")\n"
	}

	switch l.driver {
	case settings.PostgreSQL:
		return pgx
	case settings.MySQL:
		return sqlx
	}
	return ""
}

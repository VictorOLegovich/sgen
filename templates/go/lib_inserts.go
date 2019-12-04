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

//Insertion in the import section
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

func (l *libInsert) toReadOne(value, scanning string) string {
	switch l.driver {
	case settings.PostgreSQL:
		return "err := s.db.QueryRow(\n\t\tcontext.Background(),query," + value + ").Scan(\n" + scanning + "\t)"
	default:
		return ""
	}
}

func (l *libInsert) toReadAll(fillingVar, addVar, scanning string) string {
	pgxRow := "rows, err := s.db.Query(context.Background(),query.SQLString())\n\n" +
		"\tif err != nil{\n" +
		"\t\treturn nil, err\n" +
		"\t}\n\n" +
		"\tfor rows.Next(){\n" +
		"\t\terr := rows.Scan(\n" + scanning + "\t\t)\n\n" +
		"\t\tif err != nil{\n\t\t\treturn nil, err\n\t\t}\n\n" +
		"\t\t" + addVar + " = append(" + addVar + "," + fillingVar + ")\n\t}"

	switch l.driver {
	case settings.PostgreSQL:
		return pgxRow
	default:
		return ""
	}
}

func (l *libInsert) toExec(args string) string {
	pgxExec := "\t_, err := s.db.Exec(\n\t\tcontext.Background(), query,\n" + args + "\t)\n\n" +
		"\treturn err"

	switch l.driver {
	case settings.PostgreSQL:
		return pgxExec
	default:
		return ""
	}
}

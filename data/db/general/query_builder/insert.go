package query_builder

import (
	"strings"
)

type Insert struct {
	sql *strings.Builder
	*ph
}

func newInsert(sql *strings.Builder, ph *ph) *Insert {
	return &Insert{sql, ph}
}

func (i *Insert) Custom(sql string) *Insert {
	i.sql.WriteString(sql)
	return i
}

func (i *Insert) SQLString() string {
	return i.sql.String()
}

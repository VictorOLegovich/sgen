package query_builder

import (
	"strings"
)

type Delete struct {
	sql *strings.Builder
	*ph
}

func newDelete(sql *strings.Builder, ph *ph) *Delete {
	return &Delete{sql, ph}
}

func (d *Delete) Where(field, condition string) *Delete {
	checkCondition(condition)

	elems := []string{"Where ", field, " ", condition, " ", d.ph.Next()}

	for _, elem := range elems {
		d.sql.WriteString(elem)
	}

	return d
}

func (d *Delete) Custom(sql string) *Delete {
	d.sql.WriteString(sql)
	return d
}

func (d *Delete) SQLString() string {
	return d.sql.String()
}

package query_builder

const UPDATE string = `package query_builder

import (
	"strings"
)

type Update struct {
	sql *strings.Builder
	*ph
}

func newUpdate(sql *strings.Builder, ph *ph) *Update {
	return &Update{sql, ph}
}

func (u *Update) Where(field, condition string) *Update {
	checkCondition(condition)

	elems := []string{"Where ", field, " ", condition, " ", u.ph.Next()}

	for _, elem := range elems {
		u.sql.WriteString(elem)
	}

	return u
}

func (u *Update) Custom(sql string) *Update {
	u.sql.WriteString(sql)
	return u
}

func (u *Update) SQLString() string {
	return u.sql.String()
}
`

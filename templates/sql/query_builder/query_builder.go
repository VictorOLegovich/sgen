package query_builder

const QueryBuilder string = `package query_builder

import (
	"os"
	"strings"
)

type QueryBuilder struct {
	Table, Driver    string
	USet, ISet, SSet []string
	setsInit         bool
	*ph
}

const Set string = "set"

func NewQueryBuilder(Table string, Driver string) *QueryBuilder {
	return &QueryBuilder{Table: strings.ToLower(Table), Driver: Driver, ph: getPH(Driver)}
}

func (qb *QueryBuilder) InitSets(USet, ISet, SSet []string) {
	qb.USet, qb.ISet, qb.SSet, qb.setsInit = USet, ISet, SSet, true
}

func (qb *QueryBuilder) Insert() *Insert {
	if !qb.setsInit {
		println("This operation can only be used if the field sets in the database are transferred")
		os.Exit(1)
	}

	var sql strings.Builder

	elems := []string{"Insert Into ", qb.Table, " (", parameters(qb.SSet), ") Values ("}

	for i := 0; i < len(qb.ISet); i++ {
		elems = append(elems, qb.ph.Next())
		if i < len(qb.ISet)-1 {
			elems = append(elems, ", ")
		}
	}

	elems = append(elems, ")")

	for _, elem := range elems {
		sql.WriteString(elem)
	}

	return newInsert(&sql, qb.ph)
}

func (qb *QueryBuilder) Select(what string) *Select {
	var (
		sql strings.Builder
	)

	if what == Set {
		what = parameters(qb.ISet)
	}

	elems := []string{"Select ", what, "From ", strings.ToLower(qb.Table), " "}

	for _, elem := range elems {
		sql.WriteString(elem)
	}

	return newSelect(qb.Table, qb.Driver, &sql, qb.ph)
}

func (qb *QueryBuilder) Update(what string) *Update {
	var sql strings.Builder

	if !qb.setsInit {
		println("This operation can only be used if the field sets in the database are transferred")
		os.Exit(1)
	}

	elems := []string{"Update ", strings.ToLower(qb.Table), " Set"}

	if what == Set {
		for k, field := range qb.USet {
			elems = append(elems, " ", field, " = ", qb.ph.Next())
			if k < len(qb.USet)-1 {
				elems = append(elems, ", ")
			} else {
				elems = append(elems, " ")
			}
		}
	} else {
		elems = append(elems, " ", what, " = ", qb.ph.Next())
	}

	sql.WriteString(strings.Join(elems, ""))

	return newUpdate(&sql, qb.ph)
}

func (qb *QueryBuilder) Delete() *Delete {
	var sql strings.Builder

	elems := []string{"Delete From ", qb.Table, " "}

	for _, elem := range elems {
		sql.WriteString(elem)
	}

	return newDelete(&sql, qb.ph)
}
`

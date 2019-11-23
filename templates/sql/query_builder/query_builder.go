package query_builder

import (
	"strings"
)

type QueryBuilder struct {
	Table, Driver    string
	USet, CSet, SSet []string
	*ph
}

func NewQueryBuilder(Table string, USet, CSet, SSet []string, Driver string) *QueryBuilder {
	return &QueryBuilder{strings.ToLower(Table), Driver, USet, CSet, SSet, getPH(Driver)}
}

func (qb *QueryBuilder) Insert() *Insert {
	var sql strings.Builder

	elems := []string{"Insert Into ", qb.Table, " (", parameters(qb.SSet), ") Values ("}

	for i := 0; i < len(qb.CSet); i++ {
		elems = append(elems, qb.ph.Next())
		if i < len(qb.CSet)-1 {
			elems = append(elems, ", ")
		}
	}

	elems = append(elems, ")")

	for _, elem := range elems {
		sql.WriteString(elem)
	}

	return newInsert(&sql, qb.ph)
}

func (qb *QueryBuilder) Select(all bool) *Select {
	var (
		sql  strings.Builder
		what string
	)

	if all {
		what = "*"
	} else {
		what = parameters(qb.CSet)
	}

	elems := []string{"Select ", what, "From ", strings.ToLower(qb.Table), " "}

	for _, elem := range elems {
		sql.WriteString(elem)
	}

	return newSelect(qb.Table, qb.Driver, &sql, qb.ph)
}

func (qb *QueryBuilder) Update(field string) *Update {
	var sql strings.Builder

	elems := []string{"Update `", strings.ToLower(qb.Table), "` Set ", field, " = ", qb.ph.Next(), " "}

	for _, elem := range elems {
		sql.WriteString(elem)
	}

	return newUpdate(&sql, qb.ph)

}

func (qb *QueryBuilder) UpdateSeveral() *Update {
	var sql strings.Builder

	elems := []string{
		"Update `", strings.ToLower(qb.Table), "` Set"}

	for k, field := range qb.USet {
		elems = append(elems, " ", field, " = ", qb.ph.Next())
		if k < len(qb.USet)-1 {
			elems = append(elems, ", ")
		} else {
			elems = append(elems, " ")
		}
	}

	for _, elem := range elems {
		sql.WriteString(elem)
	}

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

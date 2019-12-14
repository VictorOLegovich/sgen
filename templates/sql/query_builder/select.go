package query_builder

const SELECT string = `package query_builder

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Select struct {
	table, driver string
	sql           *strings.Builder
	*ph
	sops []sop
}

func newSelect(table, driver string, sql *strings.Builder, ph *ph) *Select {
	return &Select{table, driver, sql, ph, []sop{}}
}

//////////////////////////////////////
////////// Select operations /////////
//////////////////////////////////////
func (s *Select) Where(field, condition string) *Select {
	checkCondition(condition)

	if err := s.addOp(where); err != nil {
		println(err.Error())
		os.Exit(1)
	}
	elems := []string{"Where ", field, " ", condition, " ", s.ph.Next(), " "}

	for _, elem := range elems {
		s.sql.WriteString(elem)
	}

	return s
}

func (s *Select) Limit(limitation int) *Select {
	if err := s.addOp(limit); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	elems := []string{"Limit  ", strconv.FormatInt(int64(limitation), 10), " "}

	for _, elem := range elems {
		s.sql.WriteString(elem)
	}

	return s
}

func (s *Select) OrderBy(field string, order OrderMode) *Select {
	if err := s.addOp(orderBy); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	elems := []string{"Order By ", field, " ", order.string(), " "}

	for _, elem := range elems {
		s.sql.WriteString(elem)
	}

	return s
}

func (s *Select) GroupBy(field string) *Select {
	s.sql.WriteString(" Group By ")
	s.sql.WriteString(field)

	return s
}

func (s *Select) Join(Join Join, ThisField, AttachedField, AttachedTable, Condition string) *Select {
	checkCondition(Condition)

	if err := s.addOp(join); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	sqlCollection := []string{
		Join.string(), " Join ", AttachedTable, " On ",
		s.table, ".", ThisField, Condition, AttachedTable, ".", AttachedField, " ",
	}

	for _, sqlWord := range sqlCollection {
		s.sql.WriteString(sqlWord)
	}

	return s
}

func (s *Select) Custom(sql string) *Select {
	s.sql.WriteString(sql)
	return s
}

func (s *Select) SQLString() string {
	return s.sql.String()
}

//////////////////////////////////////
////////// Sequence control //////////
//////////////////////////////////////

func (s *Select) addOp(op sop) error {
	var eBuilder strings.Builder

	eBuilder.WriteString(op.string())
	eBuilder.WriteString(" surgery is not allowed after ")

	if !s.opTolerance(op) {
		eBuilder.WriteString(sopString[s.sops[len(s.sops)-1]])
		eBuilder.WriteString("\ncompleted operations: \n")

		for key, s := range s.sops {
			eBuilder.WriteString(strconv.FormatInt(int64(key), 10))
			eBuilder.WriteString(". ")
			eBuilder.WriteString(s.string())
			eBuilder.WriteString(";\n")
		}

		return errors.New(eBuilder.String())
	}
	s.sops = append(s.sops, op)
	return nil
}

func (s *Select) opTolerance(op sop) bool {
	acceptableOps := map[sop][]sop{
		limit:   nil,
		orderBy: {limit},
		where:   {orderBy, groupBy, limit},
		join:    {orderBy, groupBy, limit, where, join},
		groupBy: {orderBy, limit},
	}

	if len(s.sops) == 0 {
		return true
	}

	lastOp := s.sops[len(s.sops)-1]

	if acceptableOps[lastOp] == nil {
		return false
	}

	for _, operation := range acceptableOps[lastOp] {
		if operation == op {
			return true
		}
	}

	return false
}

///////////////////////////////////////////////////////
///// Auxiliary types and their descriptions //////////
///////////////////////////////////////////////////////

type (
	OrderMode int
	Join      int
	sop       int
)

const (
	Desc OrderMode = iota
	Asc
)

const (
	I Join = iota
	O
	L
	R
	LO
	LI
	RO
	RI
	Empty
)

const (
	orderBy sop = iota
	groupBy
	where
	limit
	join
)

var modeString = map[OrderMode]string{
	Desc: "Desc",
	Asc:  "Asc",
}

var joinString = map[Join]string{
	I:     "Inner",
	O:     "Outer",
	L:     "Left",
	R:     "Right",
	LO:    "Left Outer",
	LI:    "Left Inner",
	RO:    "Right Outer",
	RI:    "Right Inner",
	Empty: "",
}

var sopString = map[sop]string{
	orderBy: "Order By",
	groupBy: "Group By",
	where:   "Where",
	limit:   "Limit",
	join:    "Join",
}

func (om OrderMode) string() string {
	return modeString[om]
}

func (j Join) string() string {
	return joinString[j]
}

func (s sop) string() string {
	return sopString[s]
}
`

package query_builder

import (
	"errors"
	"os"
	"strconv"
)

type ContinuedSelected struct {
	SqlString string
	perfOps   []perfOp
}

func NewContinuedSelected(sqlString string) *ContinuedSelected {
	return &ContinuedSelected{SqlString: sqlString, perfOps: make([]perfOp, 10)}
}

type perfOp int

const (
	where perfOp = iota
	limit
	orderBy
	join
)

var perfOpString = map[perfOp]string{
	where:   "Where",
	limit:   "Limit",
	join:    "Join",
	orderBy: "Order By",
}

func (cs *ContinuedSelected) Where(desired Desired) *ContinuedSelected {
	if err := cs.addOp(where); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	cs.SqlString += "Where " + desired.Field + " " + desired.Condition + " " + desired.Value + " "
	return cs
}

func (cs *ContinuedSelected) Limit(limitation int) *ContinuedSelected {
	if err := cs.addOp(limit); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	cs.SqlString += "Limit  " + strconv.Itoa(limitation) + " "
	return cs
}

func (cs *ContinuedSelected) OrderBy(order OrderBy) *ContinuedSelected {
	if err := cs.addOp(orderBy); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	cs.SqlString += "Order By " + order.OrderedField + " " + OrderMode.string()
	return cs
}

func (cs *ContinuedSelected) Join() *ContinuedSelected {
	if err := cs.addOp(where); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	cs.SqlString += "Join..."
	return cs
}

func (cs *ContinuedSelected) addOp(op perfOp) error {
	if !cs.opTolerance(op) {
		lastOp := perfOpString[cs.perfOps[len(cs.perfOps)-1]]
		return errors.New(perfOpString[op] + " surgery is not allowed after " + lastOp)
	}
	cs.perfOps = append(cs.perfOps, op)
	return nil
}

func (cs *ContinuedSelected) opTolerance(op perfOp) bool {
	acceptableOps := map[perfOp][]perfOp{
		limit:   {orderBy},
		where:   {limit, orderBy, join},
		orderBy: {limit},
		join:    {limit, orderBy, join},
	}

	if len(cs.perfOps) == 0 {
		return true
	}

	lastOp := cs.perfOps[len(cs.perfOps)-1]

	for _, operation := range acceptableOps[lastOp] {
		if operation == op {
			return true
		}
	}

	return false
}

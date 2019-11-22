package query_builder

import "errors"

type OrderBy struct {
	OrderMode
	OrderedField string
}

type OrderMode int

const (
	Desc OrderMode = iota
	//...
)

func NewOB(OrderMode OrderMode, OrderedField string) (*OrderBy, error) {
	if !fieldExist(OrderedField) {
		return nil, errors.New("a field that does not exist has been transferred to the `Order By` sql-command. ")
	}

	return &OrderBy{OrderMode, OrderedField}, nil
}

var modeToString = map[OrderMode]string{
	Desc: "Desc",
}

func (om *OrderMode) string() string {
	return modeToString[*om]
}

func fieldExist(of string) bool {
	return true
}

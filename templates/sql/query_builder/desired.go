package query_builder

import "errors"

type Desired struct {
	Field, Value, Condition string
}

func NewDesired(field, value, condition string) (*Desired, error) {
	if !fieldValid(field) {
		return nil, errors.New("the field is not valid for a 'Where' sql-condition. ")
	}

	if !conditionValid(condition) {
		return nil, errors.New("the comparison operator is not valid for a 'Where' sql-condition. ")
	}

	return &Desired{Field: field, Value: value, Condition: condition}, nil
}

func fieldValid(field string) bool {
	return true
}

func conditionValid(condition string) bool {
	conditions := []string{
		"<", ">", "<=", "=>", "=",
	}

	for _, cond := range conditions {
		if cond == condition {
			return true
		}
	}
	return false
}

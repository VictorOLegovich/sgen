package query_builder

import (
	"github.com/VictorOLegovich/sgen/collection"
	"strings"
)

type QueryBuilder struct {
	SqlString string
	*ph
}

func (qb *QueryBuilder) Select(Struct collection.Struct) *ContinuedSelected {
	return NewContinuedSelected(
		"Select " + parameters(Struct.Fields) + " From `" + strings.ToLower(Struct.Name) + "` ")
}

func (qb *QueryBuilder) Update() {

}

func parameters(parameters []collection.Field) (p string) {
	for k, parameter := range parameters {
		if k < len(parameters)-1 {
			p += parameter.Name + ", "
		} else if k == len(parameters)-1 {
			p += parameter.Name + " "
		}
	}
	return p
}

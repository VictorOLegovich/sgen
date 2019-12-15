package query_builder

const UTILS string = `package query_builder

import (
	"github.com/victorolegovich/sgen/settings"
	"os"
	"strconv"
	"strings"
)

//Utils  functions
func parameters(parameters []string) string {
	var paramBuilder strings.Builder

	for k, parameter := range parameters {
		paramBuilder.WriteString(parameter)
		if k < len(parameters)-1 {
			paramBuilder.WriteString(",")
		}
		paramBuilder.WriteString(" ")
	}
	return paramBuilder.String()
}

func checkCondition(condition string) {
	var found bool
	conditions := []string{
		"<", ">", "<=", "=>", "=", "!=",
	}

	for _, cond := range conditions {
		if cond == condition {
			found = true
		}
	}

	if !found {
		println("The wrong condition. Maybe only: [ =, !=, <=, =>, <, >")
		os.Exit(1)
	}
}

//Placeholders section
const (
	mysql     = "?"
	postgrsql = "$"
)

type ph struct {
	self   string
	inc    bool
	incNum int
}

func (ph *ph) Next() string {
	var placeholder strings.Builder

	if ph.inc {
		placeholder.WriteString(ph.self)
		placeholder.WriteString(strconv.FormatInt(int64(ph.incNum), 10))
		ph.incNum++
	} else {
		placeholder.WriteString(ph.self)
	}

	return placeholder.String()
}

func getPH(driver string) *ph {
	ph := &ph{}

	switch driver {
	case settings.MySQL:
		ph.self = mysql
		ph.inc = false
		return ph
	case settings.PostgreSQL:
		ph.self = postgrsql
		ph.inc = true
		ph.incNum = 1
		return ph
	default:
		return ph
	}
}
`

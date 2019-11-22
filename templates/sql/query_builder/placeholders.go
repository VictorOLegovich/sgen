package query_builder

import (
	"github.com/VictorOLegovich/sgen/settings"
	"strconv"
)

const (
	mysql     = "?"
	postgrsql = "$"
)

type ph struct {
	self   string
	inc    bool
	incNum int
}

func newPH(driver string) *ph {

	return getPH(driver)
}

func (ph *ph) Next() string {
	var placeholder string

	if ph.inc {
		placeholder = ph.self + strconv.Itoa(ph.incNum)
		ph.incNum++
	} else {
		placeholder = ph.self
	}

	return placeholder
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
		return ph
	default:
		return ph
	}
}

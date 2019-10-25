package types

import (
	"regexp"
)

type MapType struct {
	Key   string
	Value string
}

func IsSimpleType(typename string) bool {
	SimpleTypes := []string{
		"int8", "int16", "int32", "int64", "int",
		"uint8", "uint16", "uint32", "uint64", "uint",
		"float32", "float64", "string", "byte", "interface{}",
	}

	for _, simpleType := range SimpleTypes {
		if typename == simpleType {
			return true
		}
	}
	return false
}

func IsMap(typename string) (result bool, maptype MapType) {
	re := regexp.MustCompile("(map)[\\[0-9\\]]+[A-z]+")

	if re.MatchString(typename) {
		result = true

		for i := 4; i < len(typename); i++ {
			if typename[i] == ']' {
				maptype.Key = typename[4:i]
				maptype.Value = typename[i+1:]
			}
		}
	}

	return result, maptype
}

func IsArray(typename string) (result bool, elt string) {
	re := regexp.MustCompile("[\\[0-9\\]]+[A-z]+")

	if re.MatchString(typename) {
		result = true
		for key, sym := range typename {
			if sym == ']' {
				elt = typename[key+1:]
			}
		}
	}

	return result, elt
}

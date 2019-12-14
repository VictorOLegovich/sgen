package _go

import (
	"github.com/victorolegovich/sgen/collection"
	"github.com/victorolegovich/sgen/types"
	"regexp"
	"strings"
)

const (
	errCheck string = "\tif err != nil {\n\t\treturn err\n\t}\n\n"
	errVar   string = "var err error"
)

func err(returning string, tabs int) string {
	var tab string
	for i := 0; i < tabs; i++ {
		tab += "\t"
	}
	return tab + "if err != nil {\n" + tab + "\treturn " + returning + "\n" + tab + "}\n\n"
}

func formatTheCamelCase(s string) string {
	var builder strings.Builder
	uppercase := regexp.MustCompile("[A-Z]")

	for key, symbol := range s {
		if uppercase.MatchString(string(symbol)) {
			if key > 0 {
				if !uppercase.MatchString(string(s[key-1])) {
					builder.WriteString("_" + strings.ToLower(string(s[key])))
				} else {
					builder.WriteString(strings.ToLower(string(s[key])))
				}
			} else {
				builder.WriteString(strings.ToLower(string(s[key])))
			}
		} else {
			builder.WriteString(string(symbol))
		}
	}

	return builder.String()
}

func shortSyntaxOfCamelcase(s string) string {
	var builder strings.Builder
	uppercase := regexp.MustCompile("[A-Z]")

	for key, symbol := range s {
		if uppercase.MatchString(string(symbol)) {
			builder.WriteString(strings.ToLower(string(s[key])))
		}
	}

	return builder.String()
}

func scanningPreparation(
	varName string,
	fields []collection.Field,
	lineBreak, withId bool,
	tabs int,
) (prepared string) {

	var tabPrefix string

	for i := 0; i <= tabs; i++ {
		tabPrefix += "\t"
	}

	for _, field := range fields {
		if field.Name == "ID" && withId {
			if types.IsSimpleType(field.Type) {
				if lineBreak {
					prepared += tabPrefix + "&" + varName + "." + field.Name + ",\n"
				} else {
					prepared += "&" + varName + "." + field.Name + ","
				}
			}
			continue
		} else if field.Name == "ID" && !withId {
			continue
		}
		if types.IsSimpleType(field.Type) {
			if lineBreak {
				prepared += tabPrefix + "&" + varName + "." + field.Name + ",\n"
			} else {
				prepared += "&" + varName + "." + field.Name + ","
			}
		}
	}
	return prepared
}

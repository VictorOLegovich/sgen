package _go

import (
	"github.com/victorolegovich/sgen/collection"
	"github.com/victorolegovich/sgen/types"
	"regexp"
	"strings"
)

func parameters(fields []collection.Field) (parameters string) {
	for key, field := range fields {
		if key < len(fields) {
			parameters += field.Name + " " + field.Type + ", "
		} else {
			parameters += field.Name + " " + field.Type
		}

	}
	return parameters
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
		if types.IsSimpleType(field.Type) {
			if lineBreak {
				prepared += tabPrefix + "&" + varName + "." + field.Name + ",\n"
			} else {
				prepared += tabPrefix + "&" + varName + "." + field.Name + ","
			}
		}
	}
	return prepared
}

func hasNestedStructs(Struct collection.Struct) bool {
	for _, field := range Struct.Fields {
		if !types.IsSimpleType(field.Type) {
			return true
		}
	}
	return false
}

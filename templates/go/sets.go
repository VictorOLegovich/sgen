package _go

import (
	"github.com/victorolegovich/sgen/collection"
	"strings"
)

func (*Template) sets(Struct collection.Struct) string {
	setList := "import \"strings\"\n\n" +
		"var " + strings.ToLower(Struct.Name) + "Set = []string{\n"

	setCase := "func " + strings.ToLower(Struct.Name) + "SetCase()string{\n\t var builder strings.Builder" +
		"\tfor k, field := range " + strings.ToLower(Struct.Name) + "Set{\n" +
		"\t\tbuilder.WriteString(field)\n" +
		"\t\tif k < len(" + strings.ToLower(Struct.Name) + "Set){\n" +
		"\t\t\t builder.WriteString(\",\")\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\treturn builder.String()\n}"

	for _, field := range Struct.Fields {
		setList += "\"" + field.Name + "\","
	}

	setList += "\n}"

	return setList + setCase
}

package _go

import (
	"github.com/victorolegovich/sgen/collection"
	"github.com/victorolegovich/sgen/types"
	"strings"
)

func (t *Template) sets(Struct collection.Struct) string {
	var (
		fileContent, iSet, uSet, sSet strings.Builder
	)

	fileContent.WriteString(t.packaging(Struct.Name))

	iSet.WriteString("\nvar insertSet = []string{\n\t")
	uSet.WriteString("\nvar updateSet = []string{\n\t")
	sSet.WriteString("\nvar selectSet = []string{\n\t")

	for _, field := range Struct.Fields {
		if types.IsSimpleType(field.Type) {
			content := []string{"\"", field.Name, "\","}
			for _, c := range content {
				fc := formatTheCamelCase(c)
				if field.Name == "ID" {
					sSet.WriteString(fc)
					continue
				}
				iSet.WriteString(fc)
				uSet.WriteString(fc)
				sSet.WriteString(fc)
			}
		}
	}

	iSet.WriteString("\n}")
	uSet.WriteString("\n}")
	sSet.WriteString("\n}")

	fileContent.WriteString(iSet.String())
	fileContent.WriteString(uSet.String())
	fileContent.WriteString(sSet.String())

	return fileContent.String()
}

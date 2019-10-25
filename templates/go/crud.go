package _go

import (
	"github.com/victorolegovich/storage_generator/collection"
)

func (template *Template) crud(Struct collection.Struct) (crud string) {
	for _, crudop := range template.getcrudlist() {
		crud += crudop(Struct)
	}
	return crud
}

type crudlist []func(Struct collection.Struct) string

func (template *Template) getcrudlist() crudlist {
	return crudlist{
		template.c,
		template.r,
		template.u,
		template.d,
	}
}

func (template *Template) c(Struct collection.Struct) string {
	return "func (" + Struct.Name + "Storage *" + Struct.Name + "Storage)" +
		"Create" + Struct.Name + "(" + template.parameters(Struct.Fields) + ") " +
		" bool{\n\treturn true \n}\n\n"
}

func (template *Template) r(Struct collection.Struct) string {
	ReadOne := "func (" + Struct.Name + "Storage *" + Struct.Name + "Storage)" +
		"ReadOne" + Struct.Name + "(ID int) " + Struct.Name + "{\n\treturn " + Struct.Name + "{}\n}\n\n"

	ReadAll := "func (" + Struct.Name + "Storage *" + Struct.Name + "Storage)" +
		"Read" + Struct.Name + "List() []" + Struct.Name + "{\n\treturn []" + Struct.Name + "{}\n}\n\n"

	return ReadOne + ReadAll
}

func (template *Template) u(Struct collection.Struct) string {
	return "func (" + Struct.Name + "Storage *" + Struct.Name + "Storage)Update" + Struct.Name +
		"(Field, Value string) bool {\n\treturn true \n}\n\n"
}

func (template *Template) d(Struct collection.Struct) string {
	return "func (" + Struct.Name + "Storage *" + Struct.Name + "Storage)Delete" + Struct.Name +
		"(ID int) bool {\n\treturn true\n}\n\n"
}

func (template *Template) parameters(fields []collection.Field) (parameters string) {
	for key, field := range fields {
		if key < len(fields) {
			parameters += field.Name + " " + field.Type + ", "
		} else {
			parameters += field.Name + " " + field.Type
		}

	}
	return parameters
}

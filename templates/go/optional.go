package _go

import (
	"github.com/victorolegovich/sgen/collection"
	"strings"
)

func (t *Template) imports(Struct collection.Struct) (imports string) {
	schema := t.collection.GetRootSchema(Struct.Name)

	imported := make([]string, len(schema.Childes))

	isImported := func(imported []string, imp string) bool {
		for _, i := range imported {
			if imp == i {
				return true
			}
		}
		return false
	}

	imports += "import (\n"
	imports += t.libInsert.toImport()
	imports += "\t. \"" + t.settings.DataIA + "\"\n"
	imports += "\tqb \"" + t.settings.DatabaseIA + "/general/query_builder\"\n"

	for _, child := range schema.Childes {
		if !isImported(imported, child.StructName) {
			imports += "\t\"" + t.settings.DatabaseIA + "/" + strings.ToLower(child.StructName) + "_storage\"\n"
			imported = append(imported, child.StructName)
		}
	}

	imports += ")\n\n"
	return imports
}

func (t *Template) getByParentId(Struct collection.Struct) (withParent string) {
	schema := t.collection.GetRootSchema(Struct.Name)

	if len(schema.Parents) == 0 {
		return ""
	}

	for _, parent := range schema.Parents {
		for _, child := range t.collection.GetRootSchema(parent.StructName).Childes {
			if schema.Current == child.StructName {
				withParent += "func (s *" + schema.Current + "Storage) Get" + child.Name +
					"By" + parent.StructName + "ID(" + parent.StructName + "Id int) (" + child.Name + " " +
					child.Type + ") {\n\treturn " + child.Name + "\n}\n\n"

			}
		}

	}

	return withParent
}

func (t *Template) getWithChildes(Struct collection.Struct) string {
	var (
		storagesInit, storagesCallable, getting, adding string
	)

	schema := t.collection.GetRootSchema(Struct.Name)

	if len(schema.Childes) == 0 {
		return ""
	}

	getOne := "func (s *" + Struct.Name + "Storage)" +
		"GetOne" + Struct.Name + "WithChildes(" + Struct.Name + "ID int) " + Struct.Name + " {\n\t" +
		schema.Current + " := " + schema.Current + "Storage.ReadOne" + schema.Current + "(" + Struct.Name + "ID)\n\n"

	getAll := "func (s *" + Struct.Name + "Storage)" +
		"Get" + Struct.Name + "ListWithChildes() []" + Struct.Name + " {\n" +
		"\t" + schema.Current + "List := " + "s.ReadList()\n" +
		"\t" + schema.Current + "Buffer := make([]" + schema.Current + ", len(" + schema.Current + "List))\n\n" +
		"\tfor _, " + strings.ToLower(schema.Current) + " := range " + schema.Current + "List {" +
		"\t\t" + strings.ToLower(schema.Current) + " = " + schema.Current + "Storage.GetOne" + schema.Current +
		"WithChildes(" + strings.ToLower(schema.Current) + ".ID)\n\t\t" + schema.Current +
		"Buffer = append(" + schema.Current + "Buffer," + strings.ToLower(schema.Current) + ")\n\t}\n\n" +
		"\t return " + schema.Current + "Buffer\n}\n\n"

	for _, child := range schema.Childes {
		storageVar := strings.ToLower(child.StructName) + "Storage"
		dataVar := strings.ToLower(child.Name)

		if storagesInit == "" {
			storagesInit += "\t" + storageVar + " := " + strings.ToLower(child.StructName) + "_storage." +
				"New" + child.StructName + "Storage(nil)\n	"
		}

		getting += "\t" + strings.ToLower(child.Name) + " := " + storageVar + ".Get" + child.Name +
			"By" + schema.Current + "ID(" + schema.Current + "ID)\n"

		adding += "\t" + schema.Current + "." + child.Name + " = " + dataVar + "\n"

	}

	getOne += storagesInit + "\n" + storagesCallable + "\n" + getting + "\n" + adding +
		"\n\treturn " + schema.Current + "\n}\n\n"

	return getOne + getAll
}

func (t *Template) optionalInsert(Struct collection.Struct) string {
	schema := t.collection.GetRootSchema(Struct.Name)

	var (
		funcDecl, funcBody, storagesInit, storageCall string
	)

	funcDecl = "func (s *" + Struct.Name + ")Create(" + strings.ToLower(Struct.Name) + Struct.Name + ") error\n"

	funcBody = "{\n\tquery := s.qb.Insert()\n" +
		"\t" + t.libInsert.toExec(
		scanningPreparation(
			strings.ToLower(Struct.Name),
			Struct.Fields,
			true,
			false,
			1,
		))

	for _, child := range schema.Childes {
		storageVar := strings.ToLower(child.StructName) + "Storage"

		if storagesInit == "" {
			storagesInit += "\t" + storageVar + " := " + strings.ToLower(child.StructName) + "_storage." +
				"New" + child.StructName + "Storage(nil)\n\n"
			storageCall += "\terr := " + storageVar + ".Create(" +
				shortSyntaxOfCamelcase(Struct.Name) + "." + child.Name + ")\n\n" +
				"\t if err != nil {\n\t\treturn err\n\t}\n\n"
		}

	}

	return funcDecl + funcBody

}

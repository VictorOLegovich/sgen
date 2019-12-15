package _go

import (
	"github.com/victorolegovich/sgen/collection"
	"strings"
)

func (t *Template) imports(Struct collection.Struct) (imports string) {
	imported := make([]string, len(Struct.Childes))

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

	for _, child := range Struct.Childes {
		if !isImported(imported, child.StructName) {
			imports += "\t\"" + t.settings.DatabaseIA + "/storages/" + strings.ToLower(child.StructName) + "_storage\"\n"
			imported = append(imported, child.StructName)
		}
	}

	imports += ")\n\n"
	return imports
}

func (t *Template) getByParentId(Struct collection.Struct) (withParent string) {
	name := Struct.Name

	variable := shortSyntaxOfCamelcase(Struct.Name)
	variableDecl := variable + " := &" + Struct.Name + "{}"

	lb := len(Struct.Fields) > 4

	preparation := scanningPreparation(variable, Struct.Fields, lb, true, 1)

	if len(Struct.Parents) == 0 {
		return ""
	}

	for _, parent := range Struct.Parents {
		for _, s := range t.collection.Structs {
			if parent == s.Name {
				for _, child := range s.Childes {
					if child.StructName == name {
						var star string

						if string(child.Type[0]) != "*" {
							star = "*"
						}

						decl := "func (s *" + name + "Storage) Get" + name + "By" +
							parent + "ID(" + parent + "ID int) (" + child.Type + ", error)"

						dbCallROne := t.libInsert.toReadOne(parent+"ID", preparation, lb)
						pID := formatTheCamelCase(parent + "ID")
						query := `query := s.qb.Select(qb.Set).Where("` + pID + `","=").Limit(1).SQLString()`

						body := "{\n" +
							"\t" + errVar + "\n" +
							"\t" + variableDecl + "\n\n" +
							"\t" + query + "\n\n\t" +
							dbCallROne + "\n\n" +
							"\treturn " + star + variable + ", err\n}\n\n"

						withParent += decl + body
					}
				}
			}
		}

	}

	return withParent
}

func (t *Template) getWithChildes(Struct collection.Struct) string {
	return t.optionalOne(Struct) + t.optionalList(Struct)
}

func (t *Template) optionalOne(Struct collection.Struct) string {
	var (
		csName, storagesInit string
		getting, adding      string
		decl, body           string
		name                 = Struct.Name
		lcName               = strings.ToLower(name)
	)

	if len(Struct.Childes) == 0 {
		return ""
	}

	decl = "func (s *" + name + "Storage)One(" + name + "ID int) (*" + name + ", error)"
	body = "{\n\t" + lcName + ", err := s.one(" + name + "ID)\n" + err("nil, err", 1)

	for _, child := range Struct.Childes {
		csName = strings.ToLower(child.StructName) + "Storage"

		getting += "\tif " + lcName + "." + child.Name + ", err = " +
			"s." + csName + ".Get" + child.Name + "By" + name + "ID(" + name + "ID);err != nil{\n" +
			"\t\tprintln(err.Error())\n\t}\n\n"
	}

	body += storagesInit + "\n" + getting + "\n" + adding + "\n\treturn " + lcName + ", nil\n}\n\n"

	return decl + body
}

func (t *Template) optionalList(Struct collection.Struct) string {
	if len(Struct.Childes) == 0 {
		return ""
	}

	var (
		csvar, storagesInit string
		decl, body, getting string
		name                = Struct.Name
		lcName              = strings.ToLower(name)
		list                = lcName + "List"
	)

	decl = "func (s *" + name + "Storage)List() ([]*" + name + ", error)"
	body = "{\n\t" + list + ", err := s.list()\n" + err("nil, err", 1)

	for _, child := range Struct.Childes {
		csvar = strings.ToLower(child.StructName) + "Storage"

		getting += "\t\tif " + lcName + "." + child.Name + ", err = " +
			"s." + csvar + ".Get" + child.Name + "By" + name + "ID(" + lcName + ".ID); err != nil{\n" +
			"\t\t\tprintln(err.Error())\n\t\t}\n"
	}

	body += storagesInit + "\n" +
		"\t for _, " + lcName + " := range " + list + "{\n" +
		getting + "\t}\n\n" +
		"\treturn " + list + ", nil\n}\n\n"

	return decl + body
}

func (t *Template) optionalExec(Struct collection.Struct, operation string) string {
	var (
		funcDecl, funcBody        string
		query, exec               string
		storagesInit, storageCall string
		name, lcName              = Struct.Name, strings.ToLower(Struct.Name)
		lb, prepNil, withId       bool
	)

	if len(Struct.Fields) > 5 {
		lb = true
	}
	switch operation {
	case "Delete", "Update":
		withId = true
	case "Create":
		withId = false
	}

	funcDecl = "func (s *" + name + "Storage)" + operation + "(" + lcName + " " + name + ") error"

	preparation := scanningPreparation(lcName, Struct.Fields, lb, withId, 1)
	if operation == "Delete" {
		preparation = lcName + ".ID"
	}

	prepNil = preparation == ""

	if !prepNil {
		var qbOp string
		where := ".Where(\"ID\",\"=\")"

		switch operation {
		case "Create":
			qbOp = "Insert()"
			where = ""
		case "Update":
			qbOp = "Update(qb.Set)"
		case "Delete":
			qbOp = "Delete()"
		}

		query = "query := s.qb." + qbOp + where + ".SQLString()"
		exec = t.libInsert.toExec(preparation, lb)
	}

	funcBody = "{\n\t" + errVar + "\n\n\t" + query + "\n" + exec + "\n"

	var amp string
	var nilCheckVar string
	for _, child := range Struct.Childes {

		csName := child.StructName
		lcCsName := strings.ToLower(csName)

		storageVar := lcCsName + "Storage"

		if string(child.Type[0]) != "*" {
			amp = "&"
		}

		if amp != "" {
			nilCheckVar = "(&" + lcName + "." + child.Name + ")"
		} else {
			nilCheckVar = lcName + "." + child.Name
		}

		storageCall += "\n\tif " + nilCheckVar + " != nil{" +
			"\n\t\tif err = s." + storageVar + "." + operation + "(" + amp + lcName + "." + child.Name + ");" +
			" err != nil {\n\t\treturn err\n\t}\n\t}\n"
	}

	funcBody += storagesInit + storageCall + "\n\treturn err\n}\n\n"

	return funcDecl + funcBody

}

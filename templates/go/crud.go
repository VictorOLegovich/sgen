package _go

import (
	"github.com/victorolegovich/sgen/collection"
	"strings"
)

func (t *Template) crud(Struct collection.Struct) (crud string) {
	for _, crudop := range t.getCrudList() {
		crud += crudop(Struct)
	}
	return crud
}

type crudList []func(Struct collection.Struct) string

func (t *Template) getCrudList() crudList {
	return crudList{t.c, t.r, t.u, t.d}
}

func (t *Template) c(Struct collection.Struct) string {
	if hasNestedStructs(Struct) {
		return t.optionalInsert(Struct)
	}

	decl := "func (s *" + Struct.Name + "Storage) " + "Create(" + strings.ToLower(Struct.Name) + " " + Struct.Name + ") error"
	query := "query := s.qb.Insert().SQLString()"

	dbCall := t.libInsert.toExec(scanningPreparation(
		strings.ToLower(Struct.Name),
		Struct.Fields,
		true,
		false,
		1),
	)

	body := "{\n\t" + query + "\n\n" + dbCall + "\n}\n\n"

	return decl + body
}

func (t *Template) u(Struct collection.Struct) string {
	//Update one field
	decl := "func (s *" + Struct.Name + "Storage) Update(" +
		strings.ToLower(Struct.Name) + " " + Struct.Name + ") error"

	query := "query := s.qb.Update().Where(\"ID\", \"=\").SQLString()"

	dbCall := t.libInsert.toExec(scanningPreparation(
		strings.ToLower(Struct.Name),
		Struct.Fields,
		true,
		false,
		1),
	)

	body := "{\n\t" + query + "\n\n" + dbCall + "\n}\n\n"

	return decl + body
}

func (t *Template) d(Struct collection.Struct) string {
	decl := "func (s *" + Struct.Name + "Storage) Delete(ID int) error"

	query := `query := s.qb.Delete().Where("ID","=").SQLString()`

	dbCall := t.libInsert.toExec("ID")

	body := "{\n\t" + query + "\n\n" + dbCall + "\n}\n\n"

	return decl + body
}

func (t *Template) r(Struct collection.Struct) string {
	return t.readOne(Struct) + t.readAll(Struct)
}

func (t *Template) readOne(Struct collection.Struct) string {
	//function declaration
	declROne := "func (s *" + Struct.Name + "Storage) ReadOne(ID int) (*" + Struct.Name + ", error)"

	//variable definition
	variableROne := shortSyntaxOfCamelcase(Struct.Name)
	variableROneDecl := variableROne + " := &" + Struct.Name + "{}"

	dbCallROne := t.libInsert.toReadOne(
		"ID",
		scanningPreparation(
			variableROne,
			Struct.Fields,
			true,
			true,
			1,
		),
	)

	//query definition
	queryROne := `query := s.qb.Select(false).Where("ID","=").Limit(1).SQLString()`

	//body shaping
	bodyROne := "{\n\t" + variableROneDecl + "\n\n\t" +
		queryROne + "\n\n\t" + dbCallROne + "\n\n\treturn " + variableROne + ", err\n}\n\n"

	return declROne + bodyROne
}

func (t *Template) readAll(Struct collection.Struct) string {
	//function declaration
	declRAll := "func (s *" + Struct.Name + "Storage) ReadList() ([]" + Struct.Name + ", error)"

	//variable definition
	variableRAll := shortSyntaxOfCamelcase(Struct.Name) + "List"
	variableRAllDecl := "var " + variableRAll + " []" + Struct.Name

	variableROne := shortSyntaxOfCamelcase(Struct.Name)
	variableROneDecl := variableROne + " := " + Struct.Name + "{}"

	dbCallRAll := t.libInsert.toReadAll(
		variableROne, variableRAll,
		scanningPreparation(
			variableROne,
			Struct.Fields,
			true,
			true,
			2),
	)

	//query definition
	queryRAll := `query := s.qb.Select(false).Limit(10)`

	//body shaping
	bodyRAll := "{\n" +
		"\t" + variableRAllDecl + "\n" +
		"\t" + variableROneDecl + "\n\n" +
		"\t" + queryRAll + "\n\n" +
		"\t" + dbCallRAll + "\n\n" +
		"\treturn " + variableRAll + ", nil\n}\n\n"

	return declRAll + bodyRAll
}

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
	var (
		lb, prepNil        bool
		name               = Struct.Name
		lcName             = strings.ToLower(Struct.Name)
		query, err, dbCall string
	)

	if len(Struct.Childes) > 0 {
		return t.optionalExec(Struct, "Create")
	}

	if len(Struct.Fields) > 4 {
		lb = true
	}

	preparation := scanningPreparation(lcName, Struct.Fields, lb, false, 1)
	prepNil = preparation == ""

	decl := "func (s *" + name + "Storage) " + "Create(" + lcName + " " + name + ") error"

	if !prepNil {
		query = "query := s.qb.Insert().SQLString()"
		err = errCheck
		dbCall = t.libInsert.toExec(preparation, lb)
	}

	body := "{\n\t" + errVar + "\n\n\t" + query + "\n\n" + dbCall + "" + err + "\treturn err\n}\n\n"

	return decl + body
}

func (t *Template) u(Struct collection.Struct) string {
	if len(Struct.Childes) > 0 {
		return t.optionalExec(Struct, "Update")
	}

	var lb bool

	if len(Struct.Fields) > 5 {
		lb = true
	}

	name := Struct.Name
	lcName := strings.ToLower(name)

	decl := "func (s *" + name + "Storage) Update(" + lcName + " " + name + ") error"
	query := "query := s.qb.Update(qb.Set).Where(\"ID\", \"=\").SQLString()"

	preparation := scanningPreparation(strings.ToLower(Struct.Name), Struct.Fields, lb, true, 1)
	dbCall := t.libInsert.toExec(preparation, lb)

	body := "{\n\t" + errVar + "\n\n\t" + query + "\n\n" + dbCall + "\n\treturn err\n}\n\n"

	return decl + body
}

func (t *Template) d(Struct collection.Struct) string {
	if len(Struct.Childes) > 0 {
		return t.optionalExec(Struct, "Delete")
	}

	name := Struct.Name
	lcName := strings.ToLower(Struct.Name)

	decl := "func (s *" + name + "Storage) Delete(" + lcName + " " + name + ") error"
	query := `query := s.qb.Delete().Where("ID","=").SQLString()`
	dbCall := t.libInsert.toExec(lcName+".ID", false)
	body := "{\n\t" + errVar + "\n\t" + query + "\n\n" + dbCall + "\n\treturn err\n}\n\n"

	return decl + body
}

func (t *Template) r(Struct collection.Struct) string {
	return t.one(Struct) + t.list(Struct)
}

func (t *Template) one(Struct collection.Struct) string {
	funcName := "One"

	if len(Struct.Childes) > 0 {
		funcName = "one"
	}

	decl := "func (s *" + Struct.Name + "Storage)" + funcName + "(ID int) (*" + Struct.Name + ", error)"

	variable := shortSyntaxOfCamelcase(Struct.Name)
	variableDecl := variable + " := &" + Struct.Name + "{}"

	lb := len(Struct.Fields) > 4

	preparation := scanningPreparation(variable, Struct.Fields, lb, true, 1)

	dbCallROne := t.libInsert.toReadOne("ID", preparation, lb)

	//query definition
	query := `query := s.qb.Select(qb.Set).Where("ID","=").Limit(1).SQLString()`

	//body shaping
	body := "{\n" +
		"\t" + errVar + "\n" +
		"\t" + variableDecl + "\n\n" +
		"\t" + query + "\n\n\t" +
		dbCallROne + "\n\n" +
		"\treturn " + variable + ", err\n}\n\n"

	return decl + body
}

func (t *Template) list(Struct collection.Struct) string {
	var funcName = "List"

	if len(Struct.Childes) > 0 {
		funcName = "list"
	}

	//function declaration
	decl := "func (s *" + Struct.Name + "Storage) " + funcName + "() ([]*" + Struct.Name + ", error)"

	//variable definition
	variable := shortSyntaxOfCamelcase(Struct.Name) + "List"
	variableDecl := "var " + variable + " []*" + Struct.Name

	variableROne := shortSyntaxOfCamelcase(Struct.Name)
	variableROneDecl := variableROne + " := &" + Struct.Name + "{}"

	lb := len(Struct.Fields) > 4
	preparation := scanningPreparation(variableROne, Struct.Fields, lb, true, 2)

	dbCall := t.libInsert.toReadAll(variableROne, variable, preparation, lb)

	//query definition
	query := `query := s.qb.Select(qb.Set).Limit(10).SQLString()`

	//body shaping
	bodyRAll := "{\n" +
		"\t" + errVar + "\n" +
		"\t" + variableDecl + "\n" +
		"\t" + variableROneDecl + "\n\n" +
		"\t" + query + "\n\n" +
		"\t" + dbCall + "\n\n" +
		"\treturn " + variable + ", nil\n}\n\n"

	return decl + bodyRAll
}

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
	return crudList{
		t.c,
		t.r,
		t.u,
		t.d,
	}
}

func (t *Template) c(Struct collection.Struct) string {
	decl := "func (s *" + Struct.Name + "Storage) " + "Create(" + strings.ToLower(Struct.Name) + " " + Struct.Name + ") error"
	query := "query := s.qb.Insert()"
	body := "{\n\t" + query + "\n\tprintln(query)\n\treturn nil\n}\n\n"

	return decl + body
}

func (t *Template) r(Struct collection.Struct) string {
	//ReadOne
	declROne := "func (s *" + Struct.Name + "Storage) ReadOne(ID int) (" + Struct.Name + ", error)"
	queryROne := `query := s.qb.Select(false).Where("ID","=")`
	bodyROne := "{\n\t" + queryROne + "\n\tprintln(query)\n\treturn " + Struct.Name + "{}, nil\n}\n\n"
	ReadOne := declROne + bodyROne

	//ReadAll
	declRAll := "func (s *" + Struct.Name + "Storage) ReadList() ([]" + Struct.Name + ", error)"
	queryRAll := `query := s.qb.Select(false)`
	bodyRAll := "{\n\t" + queryRAll + "\n\tprintln(query)\n\treturn []" + Struct.Name + "{}, nil\n}\n\n"
	ReadAll := declRAll + bodyRAll

	return ReadOne + ReadAll
}

func (t *Template) u(Struct collection.Struct) string {
	//Update one field
	declUOne := "func (s *" + Struct.Name + "Storage) Update(field, value string) error"
	queryUOne := "query := s.qb.Update(field).Where(\"ID\", \"=\")"
	bodyUOne := "{\n\t" + queryUOne + "\n\tprintln(query)\n\treturn nil \n}\n\n"

	//Update a few fields
	declUSeveral := "func (s *" + Struct.Name + "Storage) UpdateSeveral(" +
		strings.ToLower(Struct.Name) + " " + Struct.Name + ") error"
	queryUSeveral := "query := s.qb.UpdateSeveral().Where(\"ID\",\"=\")"
	bodyUSeveral := "{\n\t" + queryUSeveral + "\n\tprintln(query)\n\treturn nil \n}\n\n"

	return declUOne + bodyUOne + declUSeveral + bodyUSeveral
}

func (t *Template) d(Struct collection.Struct) string {
	decl := "func (s *" + Struct.Name + "Storage) Delete(ID int) error"
	query := `query := s.qb.Delete().Where("ID","=")`
	body := "{\n\t" + query + "\n\tprintln(query)\n\treturn nil\n}\n\n"

	return decl + body
}

func (t *Template) parameters(fields []collection.Field) (parameters string) {
	for key, field := range fields {
		if key < len(fields) {
			parameters += field.Name + " " + field.Type + ", "
		} else {
			parameters += field.Name + " " + field.Type
		}

	}
	return parameters
}

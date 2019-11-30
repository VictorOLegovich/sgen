package _go

import (
	"github.com/victorolegovich/sgen/collection"
	"path/filepath"
	"strings"
)

func (t *Template) Create() (files []File) {
	storageFile := File{}
	setFile := File{}

	for _, Struct := range t.collection.Structs {
		//owner
		storageFile.Owner = Struct.Name
		setFile.Owner = Struct.Name

		//file path
		storageFile.Path = filepath.Join(t.settings.DatabaseDir, "storages", strings.ToLower(Struct.Name)+"_storage")
		setFile.Path = filepath.Join(t.settings.DatabaseDir, "storages", strings.ToLower(Struct.Name)+"_storage")

		//file name
		storageFile.Name = strings.ToLower(Struct.Name) + "_storage.go"
		setFile.Name = strings.ToLower(Struct.Name) + "_sets.go"

		//file src
		storageFile.Src = t.mainTemplate(Struct)
		setFile.Src = t.sets(Struct)

		//adding files
		files = append(files, storageFile, setFile)
	}

	return files
}

func (t *Template) mainTemplate(Struct collection.Struct) (temp string) {
	temp += t.packaging(Struct.Name)
	temp += t.imports(Struct)
	temp += t.storageType(Struct.Name)
	temp += t.newStorage(Struct.Name)
	temp += t.crud(Struct)
	temp += t.getByParentId(Struct)
	temp += t.getWithChildes(Struct)

	return temp
}

//generation of storage structure
func (*Template) storageType(strName string) string {
	return "type " + strName + "Storage struct {\n" +
		"\tDataBase string\n" +
		"\tqb *qb.QueryBuilder\n" +
		"}\n\n"
}

//generation of the function of creating a new storage
func (t *Template) newStorage(strName string) string {
	decl := "func New" + strName + "Storage(DB string) *" + strName + "Storage"

	s1, s2, s3 := "updateSet", "insertSet", "selectSet"
	driver := "\"" + t.settings.SqlDriver + "\""

	newQb := "queryBuilder := qb.NewQueryBuilder(\"" +
		strName + "\", " + s1 + ", " + s2 + ", " + s3 + ", " + driver + ")"

	return decl + "{\n\t" + newQb + "\n\treturn &" + strName + "Storage{DataBase: DB, qb: queryBuilder}\n}\n\n"
}

//package name generation
func (*Template) packaging(strName string) string {
	return "package " + strings.ToLower(strName) + "_storage\n\n"
}

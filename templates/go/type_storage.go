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
func (t *Template) storageType(strName string) string {
	return "type " + strName + "Storage struct {\n" +
		"\tdb " + t.libInsert.toType() + "\n" +
		"\tqb *qb.QueryBuilder\n" +
		"}\n\n"
}

//generation of the function of creating a new storage
func (t *Template) newStorage(strName string) string {
	decl := "func New" + strName + "Storage(db " + t.libInsert.toType() + ") *" + strName + "Storage"

	us, is, ss := "updateSet", "insertSet", "selectSet"
	driver := "\"" + t.settings.SqlDriver + "\""

	newQb := "\tqBuilder := qb.NewQueryBuilder(\"" + formatTheCamelCase(strName) + "\", " + driver + ")\n\n"
	initSets := "\t//you can opt out of using this action\n\tqBuilder.InitSets(" + us + "," + is + "," + ss + ")"

	return decl + "{\n" + newQb + initSets + "\n\n\treturn &" + strName + "Storage{db: db, qb: qBuilder}\n}\n\n"
}

//package name generation
func (*Template) packaging(strName string) string {
	return "package " + strings.ToLower(strName) + "_storage\n\n"
}

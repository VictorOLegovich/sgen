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
	temp += t.storageType(Struct)
	temp += t.newStorage(Struct)
	temp += t.crud(Struct)
	temp += t.getByParentId(Struct)
	temp += t.getWithChildes(Struct)

	return temp
}

//generation of storage structure
func (t *Template) storageType(Struct collection.Struct) string {
	typ := "type " + Struct.Name + "Storage struct {\n" +
		"\tdb " + t.libInsert.toType() + "\n" +
		"\tqb *qb.QueryBuilder\n"
	for _, child := range Struct.Childes {
		typ += "\t" + strings.ToLower(child.Name) + "Storage *" +
			strings.ToLower(child.Name) + "_storage." + child.Name + "Storage\n"
	}
	typ += "}\n\n"
	return typ
}

//generation of the function of creating a new storage
func (t *Template) newStorage(Struct collection.Struct) string {
	var childStorages string
	var transferStorages string

	for _, child := range Struct.Childes {
		lcCs := strings.ToLower(child.Name)
		childStorages += "\t" + lcCs + "Storage := " + lcCs + "_storage.New(db)\n"
		transferStorages += ", " + lcCs + "Storage: " + lcCs + "Storage"
	}

	decl := "func New(db " + t.libInsert.toType() + ") *" + Struct.Name + "Storage"

	us, is, ss := "updateSet", "insertSet", "selectSet"
	driver := "\"" + t.settings.SqlDriver + "\""

	newQb := "\tqBuilder := qb.NewQueryBuilder(\"" + formatTheCamelCase(Struct.Name) + "\", " + driver + ")\n"
	initSets := "\tqBuilder.InitSets(" + us + "," + is + "," + ss + ")\n"
	setDBName := "\tqBuilder.SetDBName(\"\")"

	body := "{\n" +
		"\t" + childStorages + "\n" +
		"\t" + newQb + initSets + setDBName + "\n\n" +
		"\treturn &" + Struct.Name + "Storage{db: db, qb: qBuilder" + transferStorages + "}\n}\n\n"

	return decl + body
}

//package name generation
func (*Template) packaging(strName string) string {
	return "package " + strings.ToLower(strName) + "_storage\n\n"
}

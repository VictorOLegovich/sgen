package _go

import (
	"github.com/victorolegovich/storage_generator/collection"
	"strings"
)

func (template *Template) Create() (files []File) {
	file := File{}
	for _, Struct := range template.collection.Structs {
		file.Owner = Struct.Name
		file.Path = template.settings.StorageDir + "/" + strings.ToLower(Struct.Name) + "_storage"
		file.Name = strings.ToLower(Struct.Name) + "_storage.go"
		file.Src = template.mainTemplate(Struct)
		files = append(files, file)
	}

	return files
}

func (template *Template) mainTemplate(Struct collection.Struct) (temp string) {
	temp += template.packaging(Struct.Name)
	temp += template.imports(Struct)
	temp += template.storagetype(Struct.Name)
	temp += template.newstorage(Struct.Name)
	temp += template.crud(Struct)
	temp += template.getByParentId(Struct)
	temp += template.getWithChildes(Struct)

	return temp
}

func (*Template) storagetype(structname string) string {
	return "type " + structname + "Storage struct {\n\tDataBase string\n}\n\n"
}

func (*Template) newstorage(structname string) string {
	return "func New" + structname + "Storage(DB string) *" + structname + "Storage {\n\treturn &" +
		structname + "Storage{DataBase: DB}\n}\n\n"
}

func (*Template) packaging(structname string) string {
	return "package " + strings.ToLower(structname) + "_storage\n\n"
}

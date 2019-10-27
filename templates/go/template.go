package _go

import (
	"github.com/victorolegovich/storage_generator/collection"
	"github.com/victorolegovich/storage_generator/settings"
)

type Template struct {
	collection collection.Collection
	settings   settings.Settings
}

type File struct {
	owner, Name, Src, Path string
}

func NewTemplate(Collection collection.Collection) *Template {
	return &Template{
		collection: Collection,
		settings: settings.Settings{
			SqlDriver: "",
			Path: settings.Path{
				ProjectDir: "/home/victor/go/src/github.com/victorolegovich/test/",
				DataDir:    "/home/victor/go/src/github.com/victorolegovich/test/data/user/user.go",
				StorageDir: "/home/victor/go/src/github.com/victorolegovich/test/storage",
			},
			ImportAliases: settings.ImportAliases{
				DataImportAlias:    "github.com/victorolegovich/test/data",
				StorageImportAlias: "github.com/victorolegovich/test/storage",
				ProjectImportAlias: "github.com/victorolegovich/test",
			},
		},
	}
}

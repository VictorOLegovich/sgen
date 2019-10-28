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

func NewTemplate(Collection collection.Collection, s settings.Settings) *Template {
	return &Template{
		collection: Collection,
		settings:   s,
	}
}

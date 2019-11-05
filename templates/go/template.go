package _go

import (
	"github.com/victorolegovich/sgen/collection"
	"github.com/victorolegovich/sgen/settings"
)

type Template struct {
	collection collection.Collection
	settings   settings.Settings
}

type File struct {
	Owner, Name, Src, Path string
}

func NewTemplate(Collection collection.Collection, s settings.Settings) *Template {
	return &Template{
		collection: Collection,
		settings:   s,
	}
}

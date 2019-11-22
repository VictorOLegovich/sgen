package query_builder

type StructDescription struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name, Type string
}

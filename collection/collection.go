package collection

type Collection struct {
	DataPackage string
	Structs     []Struct
}

type RootSchema struct {
	Parents []string
	Childes []RootObject
}

type RootObject struct {
	StructName string
	Field
}

type Struct struct {
	Name        string
	Fields      []Field
	Complicated map[string]Complicated
	RootSchema
}

type Field struct {
	Name, Type string
}

type Complicated int

const (
	ComplicatedMap Complicated = iota
	ComplicatedArray
	ComplicatedInterface
	Empty
)

var complicatedToString = map[Complicated]string{
	ComplicatedMap:       "ComplicatedMap",
	ComplicatedArray:     "ComplicatedArray",
	ComplicatedInterface: "ComplicatedInterface",
}

func (c Complicated) String() string {
	return complicatedToString[c]
}

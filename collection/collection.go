package collection

type Collection struct {
	DataPackage string
	Structs     []Struct
	RootSchemas []RootSchema
}

type RootSchema struct {
	Current string
	Parents []RootObject
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
}

type Field struct {
	Name, Type string
}

func (collection *Collection) CompletingRootSchemas() {
	for _, schema := range collection.RootSchemas {
		for _, child := range schema.Childes {
			for key, s := range collection.RootSchemas {
				if s.Current == child.StructName {
					var rObj RootObject

					if len(schema.Parents) == 0 {
						rObj = RootObject{
							StructName: schema.Current,
							Field:      Field{},
						}
					} else {

						rObj.StructName = schema.Current

						for _, parent := range schema.Parents {
							for _, child := range collection.GetRootSchema(parent.StructName).Childes {
								if child.StructName == schema.Current {
									rObj.Field = child.Field
								}
							}
						}
					}
					if !collection.parentExist(rObj) {
						collection.RootSchemas[key].Parents =
							append(collection.RootSchemas[key].Parents, rObj)
					}
				}
			}
		}
	}
}

func (collection *Collection) parentExist(object RootObject) bool {
	for _, rs := range collection.RootSchemas {
		for _, parent := range rs.Parents {
			if parent.StructName == object.StructName {
				return true
			}
		}
	}
	return false
}

func (collection *Collection) GetRootSchema(root string) RootSchema {
	for _, schema := range collection.RootSchemas {
		if schema.Current == root {
			return schema
		}
	}
	return RootSchema{}
}

func (collection *Collection) GetStruct(structname string) Struct {
	for _, Struct := range collection.Structs {
		if Struct.Name == structname {
			return Struct
		}
	}
	return Struct{}
}

func (Struct *Struct) GetField(fieldname string) Field {
	for _, Field := range Struct.Fields {
		if Field.Name == fieldname {
			return Field
		}
	}
	return Field{}
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

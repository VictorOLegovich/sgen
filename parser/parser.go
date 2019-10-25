package parser

import (
	c "github.com/victorolegovich/storage_generator/collection"
	"go/ast"
	"go/parser"
	"go/token"
)

func Parse(filename string, collection *c.Collection) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	//_ = ast.Print(fset, file)

	inspecting(file, collection)

	collection.CompletingRootSchemas()

	return nil
}

func inspecting(file *ast.File, collection *c.Collection) {

	ast.Inspect(file, func(node ast.Node) bool {
		var Struct c.Struct
		var RootSchema c.RootSchema

		switch n := node.(type) {
		case *ast.File:
			collection.DataPackage = n.Name.Name
		case *ast.TypeSpec:

			switch s := n.Type.(type) {

			case *ast.StructType:
				Struct.Name, RootSchema.Current = n.Name.Name, n.Name.Name
				Struct.Fields, RootSchema.Childes, Struct.Complicated = fillFlds(s)

				collection.Structs = append(collection.Structs, Struct)
				collection.RootSchemas = append(collection.RootSchemas, RootSchema)
			}
		}

		return true
	})
}

//fields filling
func fillFlds(s *ast.StructType) (Fields []c.Field, Childes []c.RootObject, Complicated map[string]c.Complicated) {
	var Field = c.Field{}
	var Child = c.RootObject{}

	Complicated = map[string]c.Complicated{}

	for _, field := range s.Fields.List {
		//не до конца уверен, что так правильно, но это работает.
		//так мы получаем безымянные поля(то есть структуры и интерфейсы)
		if len(field.Names) == 0 {

			switch ftype := field.Type.(type) {

			case *ast.Ident:
				Child.StructName, Child.Type, Child.Name = ftype.Name, ftype.Name, ftype.Name
				Childes = append(Childes, Child)
				Field.Name, Field.Type = ftype.Name, ftype.Name
			}
		}

		for _, ident := range field.Names {

			Field.Name = ident.Name

			comp, child := defineFieldType(field, &Field)

			if comp != c.Empty {
				Complicated[Field.Name] = comp
			}
			if child.StructName != "" {
				child.Field = Field
				Childes = append(Childes, child)
			}
		}

		Fields = append(Fields, Field)
	}

	return Fields, Childes, Complicated
}

func defineFieldType(field *ast.Field, Field *c.Field) (complicated c.Complicated, child c.RootObject) {
	complicated = c.Empty

	switch ftype := field.Type.(type) {

	//Простой тип
	case *ast.Ident:
		Field.Type = ftype.Name

	case *ast.InterfaceType:
		if len(ftype.Methods.List) == 0 {
			Field.Type = "interface{}"
		} else {
			complicated = c.ComplicatedInterface
		}

	//Мапа
	case *ast.MapType:

		switch key := ftype.Key.(type) {

		case *ast.Ident:
			Field.Type = "map[" + key.Name + "]"
			switch typespec := key.Obj.Decl.(type) {
			case *ast.TypeSpec:
				switch typespec.Type.(type) {
				case *ast.StructType:
					complicated = c.ComplicatedMap
				}
			}
		}

		switch value := ftype.Value.(type) {

		case *ast.Ident:
			Field.Type += value.Name

			switch decl := value.Obj.Decl.(type) {
			case *ast.TypeSpec:
				switch decl.Type.(type) {
				case *ast.StructType:
					child.Type = value.Name
				}
			}

		//Это усложнённые типы, такое мы обрабатывать не будем.
		case *ast.MapType:
			complicated = c.ComplicatedMap

		case *ast.ArrayType:
			complicated = c.ComplicatedMap
		}

	//Массив
	case *ast.ArrayType:
		var arrayLen string

		switch arrlen := ftype.Len.(type) {
		case *ast.BasicLit:
			arrayLen = arrlen.Value
		}

		switch arrname := ftype.Elt.(type) {

		case *ast.Ident:
			if arrayLen != "" {
				Field.Type = "[" + arrayLen + "]" + arrname.Name
			} else {
				Field.Type = "[]" + arrname.Name
			}
			switch decl := arrname.Obj.Decl.(type) {
			case *ast.TypeSpec:
				switch decl.Type.(type) {
				case *ast.StructType:
					child.StructName = arrname.Name
				}
			}

		//Это усложнённые типы, такое мы обрабатывать не будем.
		case *ast.MapType:
			complicated = c.ComplicatedArray

		case *ast.ArrayType:
			complicated = c.ComplicatedArray
		}

	//Скорее всего, это импортируемый тип.
	case *ast.SelectorExpr:
		switch x := ftype.X.(type) {

		case *ast.Ident:
			Field.Type = x.Name + "." + ftype.Sel.Name
		}

	case *ast.StarExpr:
		switch x := ftype.X.(type) {

		case *ast.SelectorExpr:
			switch x2 := x.X.(type) {
			case *ast.Ident:
				Field.Type = "*" + x2.Name + "." + x.Sel.Name
			}
		}
	}

	return complicated, child
}

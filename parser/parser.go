package parser

import (
	"fmt"
	c "github.com/victorolegovich/storage_generator/collection"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
)

func Parse(filename string, collection *c.Collection) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	//_ = ast.Print(fset, file)

	//если нет поля ID, добавляем его в файл
	if pos := inspecting(file, collection); pos != 0 {
		addToFile(filename, pos)
	}

	collection.CompletingRootSchemas()

	return nil
}

func inspecting(file *ast.File, collection *c.Collection) (pos int) {

	ast.Inspect(file, func(node ast.Node) bool {
		var Struct c.Struct
		var RootSchema c.RootSchema

		switch n := node.(type) {
		case *ast.File:
			collection.DataPackage = n.Name.Name
		case *ast.TypeSpec:

			switch s := n.Type.(type) {

			case *ast.StructType:
				if !hasID(s.Fields) {
					pos = int(s.Struct)
				}

				Struct.Name, RootSchema.Current = n.Name.Name, n.Name.Name
				Struct.Fields, RootSchema.Childes, Struct.Complicated = fillFlds(s)

				collection.Structs = append(collection.Structs, Struct)
				collection.RootSchemas = append(collection.RootSchemas, RootSchema)
			}
		}

		return true
	})
	return pos
}

//fields filling
func fillFlds(s *ast.StructType) (Fields []c.Field, Childes []c.RootObject, Complicated map[string]c.Complicated) {
	var Field = c.Field{}
	var Child = c.RootObject{}

	Complicated = map[string]c.Complicated{}

	for _, field := range s.Fields.List {
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

		//Усложнился вложенной структурой
		if ftype.Obj != nil {
			switch decl := ftype.Obj.Decl.(type) {

			case *ast.TypeSpec:
				switch decl.Type.(type) {
				case *ast.StructType:
					Field.Type = decl.Name.Name
					child.StructName = decl.Name.Name
				}
			}
		}

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
			if arrname.Obj != nil {
				switch decl := arrname.Obj.Decl.(type) {
				case *ast.TypeSpec:
					switch decl.Type.(type) {
					case *ast.StructType:
						child.StructName = arrname.Name
					}
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

func hasID(fields *ast.FieldList) bool {
	for _, field := range fields.List {
		for _, ident := range field.Names {
			if ident.Name == "ID" {
				return true
			}
		}
	}
	return false
}

func addToFile(filename string, pos int) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	info, _ := file.Stat()
	data := make([]byte, info.Size())

	for {
		_, err := file.Read(data)

		if err == io.EOF {
			break
		}
	}
	_ = file.Close()

	adding := "\n\tID int\n"
	writing := string(data[:pos+7]) + adding + string(data[pos+len(adding)-1:])

	file, _ = os.OpenFile(filename, os.O_WRONLY, os.ModeAppend)
	_, _ = file.WriteString(writing)

	_ = file.Close()
}

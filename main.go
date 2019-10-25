package main

import (
	c "github.com/victorolegovich/storage_generator/collection"
	p "github.com/victorolegovich/storage_generator/parser"
	_go "github.com/victorolegovich/storage_generator/templates/go"
	"github.com/victorolegovich/storage_generator/validator"
	"os"
)

func main() {
	collection := &c.Collection{}

	_ = p.Parse("/home/victor/go/src/github.com/victorolegovich/test/data/user/user.go", collection)

	if err := validator.StructsValidation(collection.Structs); err != nil {
		println(err.Error())
		return
	}

	template := _go.NewTemplate(*collection)

	for _, file := range template.Create() {
		if mkerr := os.Mkdir(file.Path, os.ModePerm); mkerr != nil {
			//	println("mkdir error :  " + mkerr.Error())
		}

		if f, err := os.Create(file.Path + "/" + file.Name); err == nil {
			_, _ = f.Write([]byte(file.Src))
		} else {
			//println(err.Error())
		}
	}
}

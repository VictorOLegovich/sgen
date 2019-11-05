package main

import "github.com/victorolegovich/storage_generator/generator"

func main() {
	g := generator.Generator{}
	if err := g.Generate(); err != nil {
		println(err.Error())
	}

}

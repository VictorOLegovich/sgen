package main

import (
	"github.com/victorolegovich/sgen/settings"
	"github.com/victorolegovich/sgen/templates/sql/query_builder"
)

func main() {
	//if err := sgen_cli.Run(); err != nil {
	//	println(err.Error())
	//}

	USet := []string{"Man"}
	CSet := []string{"Man", "Skin"}
	SSet := []string{"Man", "Skin", "Sex"}

	builder := query_builder.NewQueryBuilder("user", USet, CSet, SSet, settings.PostgreSQL)

	j := query_builder.RI
	f1 := "Man"
	at := "Employee"
	f2 := "ID"

	println(builder.Select(false).Join(j, f1, f2, at, "=").Where("Man", "=").SQLString())
}

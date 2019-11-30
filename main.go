package main

import "github.com/victorolegovich/sgen/sgen_cli"

func main() {
	if err := sgen_cli.Run(); err != nil {
		println(err.Error())
	}
}

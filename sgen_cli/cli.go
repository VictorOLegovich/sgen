package sgen_cli

import (
	"errors"
	"github.com/urfave/cli"
	"github.com/victorolegovich/sgen/generator"
	"os"
)

func Run() error {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"gen", "-g", "g"},
			Usage:   "main generate command",
			Action: func(c *cli.Context) error {
				if _, err := os.Stat(c.Args().First()); os.IsNotExist(err) {
					return errors.New("File not exist: " + c.Args().First())
				}
				gen := generator.Generator{}

				if err := gen.Generate(c.Args().First()); err != nil {
					return err
				}

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}

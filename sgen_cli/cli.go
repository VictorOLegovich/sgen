package sgen_cli

import (
	"errors"
	"github.com/urfave/cli"
	"github.com/victorolegovich/sgen/generator"
	"github.com/victorolegovich/sgen/settings"
	"os"
	"path/filepath"
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
		{
			Name:    "gs",
			Aliases: []string{"settings"},
			Usage:   "getting settings file",
			Action: func(c *cli.Context) error {
				targetDir, err := filepath.Abs(c.Args().First())
				if err != nil {
					return err
				}

				if _, err := os.Stat(c.Args().First()); os.IsNotExist(err) {
					return errors.New("Dir is not exist: " + c.Args().First())
				}

				file, err := os.Create(filepath.Join(targetDir, "settings.json"))
				if err != nil {
					return err
				}

				_, err = file.WriteString(settings.SettingsSRC)
				return err
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}

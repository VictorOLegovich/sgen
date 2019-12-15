package sgen_cli

import (
	"errors"
	"github.com/urfave/cli"
	"github.com/victorolegovich/sgen/generator"
	"github.com/victorolegovich/sgen/settings"
	"io/ioutil"
	"os"
	"path/filepath"
)

const settingsFile string = "sgen.json"

func Run() error {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"gen", "-g", "g"},
			Usage:   "main generate command",
			Action: func(c *cli.Context) error {
				if _, err := os.Stat(c.Args().First()); os.IsNotExist(err) {
					if files, err := ioutil.ReadDir("."); err == nil {
						for _, file := range files {
							if file.Name() == settingsFile {
								path, _ := filepath.Abs(".")
								if err := generator.Generate(filepath.Join(path, settingsFile)); err != nil {
									return err
								} else {
									return nil
								}
							}
						}
					}
					return errors.New("File not exist: " + c.Args().First())
				}

				if err := generator.Generate(c.Args().First()); err != nil {
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
				targetDir, _ := filepath.Abs(".")

				file, err := os.Create(filepath.Join(targetDir, settingsFile))
				if err != nil {
					return err
				}

				_, err = file.WriteString(settings.SRC)
				return err
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}

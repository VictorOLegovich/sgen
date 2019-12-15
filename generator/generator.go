package generator

import (
	"errors"
	"github.com/victorolegovich/sgen/collection"
	"github.com/victorolegovich/sgen/file_manager"
	"github.com/victorolegovich/sgen/parser"
	"github.com/victorolegovich/sgen/settings"
	_go "github.com/victorolegovich/sgen/templates/go"
)

// The process of generating storage files is performed in several stages:
//     1. Receiving a file with settings;
//     2. Parsing the DataDir directory;
//     3. Template depending on the settings;
//     4. Placement of files received during the template process
func Generate(file string) error {
	//Receiving settings and initializing them
	sett, err := settings.New(file)
	if err != nil {
		return err
	}

	//Filling in the required generator fields

	//Collection initialization
	c := &collection.Collection{}

	//Parsing
	if err := parser.Parse(sett.Path.DataDir, c); err != nil {
		return perr("parsing", err)
	}

	//Template generation
	template := _go.NewTemplate(*c, sett)

	//Deploying files
	depl := file_manager.NewFileManger(sett, template.Create())
	if err := depl.Deploy(); err != nil {
		return perr("deploying", err)
	}

	return nil
}

//Generation process error
func perr(section string, err error) error {
	return errors.New(section + "section of generating process: \n\t" + err.Error() + "\n")
}

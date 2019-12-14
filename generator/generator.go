package generator

import (
	"errors"
	"github.com/victorolegovich/sgen/collection"
	"github.com/victorolegovich/sgen/file_manager"
	"github.com/victorolegovich/sgen/parser"
	"github.com/victorolegovich/sgen/settings"
	_go "github.com/victorolegovich/sgen/templates/go"
	"github.com/victorolegovich/sgen/validator"
)

type Generator struct {
	settings settings.Settings
}

// The process of generating storage files is performed in several stages:
//     1. Receiving a file with settings;
//     2. Parsing the DataDir directory;
//     3. Validation of the collection obtained by parsing
//     4. Template depending on the settings;
//     5. Placement of files received during the template process
func (gen *Generator) Generate(file string) error {
	//Receiving settings and initializing them
	s, err := settings.New(file)
	if err != nil {
		return err
	}

	//Filling in the required generator fields
	gen.settings = s

	//Collection initialization
	c := &collection.Collection{}

	//Parsing
	if err := parser.Parse(gen.settings.Path.DataDir, c); err != nil {
		return perr("parsing", err)
	}

	//Validating
	if err := validator.StructsValidation(c.Structs); err != nil {
		return perr("validating", err)
	}

	//Template generation
	template := _go.NewTemplate(*c, gen.settings)

	//Deploying files
	depl := file_manager.NewFileManger(s, template.Create())
	if err := depl.Deploy(); err != nil {
		return perr("deploying", err)
	}

	return nil
}

//Generation process error
func perr(section string, err error) error {
	return errors.New(section + "section of generating process: \n\t" + err.Error() + "\n")
}

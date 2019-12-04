package generator

import (
	"errors"
	"github.com/victorolegovich/sgen/collection"
	"github.com/victorolegovich/sgen/file_manager"
	"github.com/victorolegovich/sgen/parser"
	reg "github.com/victorolegovich/sgen/register"
	"github.com/victorolegovich/sgen/settings"
	_go "github.com/victorolegovich/sgen/templates/go"
	"github.com/victorolegovich/sgen/validator"
	"os"
)

type Generator struct {
	settings settings.Settings
}

// The process of generating storage files is performed in several stages:
//     1. Receiving a file with settings;
//     2. Initialization of the registry, if auto-delete is set as true in the settings;
//     3. Parsing the DataDir directory;
//     4. Validation of the collection obtained by parsing
//     5. Template depending on the settings;
//     6. Placement of files received during the template process
//     7. Adding/deleting data to the registry
//     8. Removal of storages that do not have structures in this package
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

	//Register object initialization
	rObj := &reg.RegObject{}
	rObj.Entistor = map[string]string{}

	//Register initialization
	r, err := reg.NewRegister()
	if err != nil {
		return perr("register", err)
	}

	//Parsing
	if err := parser.Parse(gen.settings.Path.DataDir, c); err != nil {
		return perr("parsing", err)
	} else {
		//Adding a package name to the registry object
		rObj.Package = c.DataPackage
	}

	//Validating
	if err := validator.StructsValidation(c.Structs); err != nil {
		return perr("validating", err)
	}

	//Template generation
	template := _go.NewTemplate(*c, gen.settings)

	//Deploying files
	depl := file_manager.NewFileManger(s, template.Create(), rObj)
	if err := depl.Deploy(); err != nil {
		return perr("deploying", err)
	}

	//Adding an object to the register
	r.AddObject(*rObj)

	//Saving register state
	if err := r.Save(); err != nil {
		return perr("registration", err)
	}

	if gen.settings.AutoDelete {
		for _, del := range r.Deleted {
			if err := os.RemoveAll(del); err != nil {
				println("-------------------------------------")
				println("             Removing error:")
				println("             " + err.Error())
				println("-------------------------------------")
			}
		}
	}

	return nil
}

//Generation process error
func perr(section string, err error) error {
	return errors.New(section + "section of generating process: \n\t" + err.Error() + "\n")
}

package generator

import (
	"errors"
	"github.com/victorolegovich/storage_generator/collection"
	reg "github.com/victorolegovich/storage_generator/generator/register"
	"github.com/victorolegovich/storage_generator/parser"
	"github.com/victorolegovich/storage_generator/settings"
	_go "github.com/victorolegovich/storage_generator/templates/go"
	"github.com/victorolegovich/storage_generator/validator"
	"io/ioutil"
	"os"
)

type genSection int

const (
	parsing genSection = iota
	validating
	templating
	register
)

var gsToString = map[genSection]string{
	parsing:    "parsing",
	validating: "validating",
	templating: "templating",
	register:   "register",
}

func (gs genSection) string() string {
	return gsToString[gs]
}

type processError map[genSection]error

type Generator struct {
	settings settings.Settings
	processError
}

func (gen *Generator) Generate() error {
	c := &collection.Collection{}
	r := reg.NewRegister()
	rObj := reg.NewRegObject()

	if err := parser.Parse(gen.settings.Path.DataDir, c); err != nil {
		gen.processError[parsing] = err
	} else {
		rObj.Name = c.DataPackage
		files, _ := ioutil.ReadDir(gen.settings.DataDir)
		for _, file := range files {
			if err = rObj.AddToDataFilesState(gen.settings.DataDir + "/" + file.Name()); err != nil {
				gen.processError[register].Error()
			}
		}
	}

	if err := validator.StructsValidation(c.Structs); err != nil {
		gen.processError[validating] = err
	}

	template := _go.NewTemplate(*c, gen.settings)

	for _, file := range template.Create() {
		if _, err := os.Stat(file.Path); os.IsNotExist(err) {
			if err = os.Mkdir(file.Path, os.ModePerm); err != nil {
				gen.processError[templating] = err
			}
		}

		if f, err := os.Create(file.Path + "/" + file.Name); err == nil {

			if err = rObj.AddToStorageFilesState(file.Path + "/" + file.Name); err != nil {
				gen.processError[register] = err
			}

			if _, err = f.Write([]byte(file.Src)); err != nil {
				gen.processError[templating] = err
			}
		} else {
			gen.processError[templating] = err
		}
	}

	if err := r.AddObject(*rObj); err != nil {

	}

	return errorConverting(gen.processError)
}

func errorConverting(pErr processError) error {
	var errorText string

	for section, err := range pErr {
		errorText += section.string() + " section of generating: \n\t" + err.Error() + "\n"
	}

	if errorText != "" {
		return errors.New(errorText)
	}

	return nil
}

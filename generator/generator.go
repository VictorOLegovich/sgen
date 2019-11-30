package generator

import (
	"errors"
	"fmt"
	"github.com/victorolegovich/sgen/collection"
	"github.com/victorolegovich/sgen/file_manager"
	"github.com/victorolegovich/sgen/parser"
	reg "github.com/victorolegovich/sgen/register"
	"github.com/victorolegovich/sgen/settings"
	_go "github.com/victorolegovich/sgen/templates/go"
	"github.com/victorolegovich/sgen/validator"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type genSection int

const (
	parsing genSection = iota
	validating
	deploying
	register
)

var gsToString = map[genSection]string{
	parsing:    "parsing",
	validating: "validating",
	deploying:  "deploying",
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

func (gen *Generator) Generate(file string) error {
	s, e := settings.New(file)
	if e != nil {
		return e
	}

	gen.settings = s
	gen.processError = map[genSection]error{}
	c := &collection.Collection{}
	r, e := reg.NewRegister()
	if e != nil {
		gen.processError[register] = e
	}
	rObj := &reg.RegObject{}
	rObj.Entistor = map[string]string{}

	if err := parser.Parse(gen.settings.Path.DataDir, c); err != nil {
		gen.processError[parsing] = err
	} else {
		rObj.Package = c.DataPackage
	}

	if err := validator.StructsValidation(c.Structs); err != nil {
		gen.processError[validating] = err
	}

	template := _go.NewTemplate(*c, gen.settings)

	depl := file_manager.NewFileManger(s, template.Create(), rObj)

	if err := depl.Deploy(); err != nil {
		gen.processError[deploying] = err
	}

	r.AddObject(*rObj)

	if err := r.Save(); err != nil {
		gen.processError[register] = err
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

	return errorConverting(gen.processError)
}

func errorConverting(pErr processError) error {
	var errorText string

	for section, err := range pErr {
		errorText += section.string() + " section of generating process: \n\t" + err.Error() + "\n"
	}

	if errorText != "" {
		return errors.New(errorText)
	}

	return nil
}

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

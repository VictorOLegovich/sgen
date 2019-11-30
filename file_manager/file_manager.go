package file_manager

import (
	"fmt"
	"github.com/victorolegovich/sgen/register"
	"github.com/victorolegovich/sgen/settings"
	_go "github.com/victorolegovich/sgen/templates/go"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileManager struct {
	settings settings.Settings
	files    []_go.File
	ro       *register.RegObject
}

func NewFileManger(settings settings.Settings, files []_go.File, ro *register.RegObject) *FileManager {
	return &FileManager{settings, files, ro}
}

func (fm *FileManager) Deploy() error {
	if err := fm.createBaseDirectories(); err != nil {
		return err
	}

	if err := fm.moveTheAuxiliaryModules(); err != nil {
		return err
	}

	if err := fm.createFiles(); err != nil {
		return err
	}

	return nil
}

func (fm *FileManager) createBaseDirectories() error {
	general, _ := filepath.Abs(fm.settings.DatabaseDir + "/general")
	if err := os.Mkdir(general, os.ModePerm); err != nil && os.IsNotExist(err) {
		return err
	}

	storages, _ := filepath.Abs(fm.settings.DatabaseDir + "/storages")
	if err := os.Mkdir(storages, os.ModePerm); err != nil && os.IsNotExist(err) {
		return err
	}

	return nil
}

func (fm *FileManager) moveTheAuxiliaryModules() error {
	if err := fm.moveQB(); err != nil {
		return err
	}

	return nil
}

func (fm *FileManager) moveQB() error {
	srcDir, _ := filepath.Abs("../templates/sql/query_builder")

	dstDir, _ := filepath.Abs(filepath.Join(fm.settings.DatabaseDir, "general", "query_builder"))
	if _, err := os.Stat(dstDir); os.IsExist(err) {
		println("Уже есть qb")
		return nil
	}

	return CopyDir(srcDir, dstDir)
}

func (fm *FileManager) createFiles() error {
	for _, file := range fm.files {
		if _, err := os.Stat(file.Path); os.IsNotExist(err) {
			if err = os.Mkdir(file.Path, os.ModePerm); err != nil {
				return err
			}
		}

		if f, err := os.Create(filepath.Join(file.Path, file.Name)); err == nil {
			fm.ro.Entistor[file.Owner] = file.Path

			if _, err = f.Write([]byte(file.Src)); err != nil {
				return err
			}
		} else {
			return err
		}
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

func Delete(del string) error {
	return os.RemoveAll(del)
}

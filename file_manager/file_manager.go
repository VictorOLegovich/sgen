package file_manager

import (
	"errors"
	"fmt"
	"github.com/victorolegovich/sgen/settings"
	_go "github.com/victorolegovich/sgen/templates/go"
	"github.com/victorolegovich/sgen/templates/go/general"
	"github.com/victorolegovich/sgen/templates/sql/query_builder"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Scope int

const (
	File Scope = iota
	Decl
)

type FileManager struct {
	settings settings.Settings
	files    []_go.File
}

func NewFileManger(settings settings.Settings, files []_go.File) *FileManager {
	return &FileManager{settings, files}
}

func (fm *FileManager) Deploy() error {
	if err := fm.createBaseDirectories(); err != nil {
		return err
	}

	if err := fm.moveModules(); err != nil {
		return err
	}

	if err := fm.createFiles(); err != nil {
		return err
	}

	return nil
}

func (fm *FileManager) createBaseDirectories() error {
	genDir, _ := filepath.Abs(fm.settings.DatabaseDir + "/general")
	storages, _ := filepath.Abs(fm.settings.DatabaseDir + "/storages")
	database, _ := filepath.Abs(fm.settings.DatabaseDir + "/general/db")
	queryBuilder, _ := filepath.Abs(fm.settings.DatabaseDir + "/general/query_builder")

	dirs := []string{genDir, storages, database, queryBuilder}
	for _, dir := range dirs {
		if err := os.Mkdir(dir, os.ModePerm); err != nil && os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

func (fm *FileManager) moveModules() error {
	if err := fm.moveQB(); err != nil {
		return err
	}

	if err := fm.moveDB(); err != nil {
		return err
	}

	return nil
}

func (fm *FileManager) moveQB() error {
	queryBuilder, _ := filepath.Abs(fm.settings.DatabaseDir + "/general/query_builder")

	for _, file := range query_builder.Files() {
		f, err := os.Create(filepath.Join(queryBuilder, file.Name))
		if err != nil {
			return err
		}

		if _, err = f.WriteString(file.Src); err != nil {
			return err
		}
	}

	return nil
}

func (fm *FileManager) moveDB() error {
	dst, _ := filepath.Abs(fm.settings.DatabaseDir + "/general/db/db.go")

	file, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = file.WriteString(general.Src(fm.settings.SqlDriver))

	return err
}

func (fm *FileManager) createFiles() error {
	for _, file := range fm.files {
		if _, err := os.Stat(file.Path); os.IsNotExist(err) {
			if err = os.Mkdir(file.Path, os.ModePerm); err != nil {
				return err
			}
		}

		if f, err := os.Create(filepath.Join(file.Path, file.Name)); err == nil {
			if _, err = f.Write([]byte(file.Src)); err != nil {
				return err
			}

			cmd := exec.Command("go", "fmt", file.Path)
			if err := cmd.Run(); err != nil {
				print("Fmt can't be called")
			}

		} else {
			return err
		}
	}

	return nil
}

func AddToFile(filename, exp, needle string, scope Scope) error {
	var (
		data    []string
		builder strings.Builder
		ex      = regexp.MustCompile(exp)
	)

	data = strings.Split(read(filename), "\n")

	if scope != File {

		positions := positions(data, ex)

		if len(positions) == 0 {
			return errors.New("No space was found in this file ")
		}

		for i, s := range data {
			if i > 0 {
				builder.WriteString("\n")
			}
			if hasPos(positions, i) {
				builder.WriteString(needle + "\n")
			}

			builder.WriteString(s)

		}

		write(filename, builder.String())
		return format(filename)
	}

	for k, s := range data {
		if k > 0 {
			builder.WriteString("\n")
		}

		if ex.MatchString(s) {
			builder.WriteString(needle)
		}

		builder.WriteString(s)
	}

	write(filename, builder.String())
	return format(filename)
}

func hasPos(positions []int, current int) bool {
	for _, position := range positions {
		if current == position {
			return true
		}
	}

	return false
}

func positions(data []string, exp *regexp.Regexp) (positions []int) {
	for i, s := range data {
		if exp.MatchString(s) {
			positions = append(positions, i+1)
		}
	}

	return positions
}

func read(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	info, _ := file.Stat()
	data := make([]byte, info.Size())

	for {
		_, err := file.Read(data)

		if err == io.EOF {
			break
		}
	}
	_ = file.Close()

	return string(data)
}

func write(filename, content string) {
	file, _ := os.OpenFile(filename, os.O_WRONLY, os.ModeAppend)
	_, _ = file.WriteString(content)

	_ = file.Close()
}

func format(file string) error {
	cmd := exec.Command("go", "fmt", file)
	return cmd.Run()
}

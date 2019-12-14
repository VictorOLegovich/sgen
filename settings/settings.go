package settings

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const SRC = `{
  "path": {
    "project_dir": "",
    "data_dir": "",
    "database_dir": ""
  },
  "sql_driver": ""
}`

const (
	MySQL      = "mysql"
	PostgreSQL = "postgresql"
)

type (
	Settings struct {
		Path `json:"path"`
		ImportAliases
		SqlDriver string `json:"sql_driver"`
	}
	Path struct {
		ProjectDir  string `json:"project_dir"`
		DataDir     string `json:"data_dir"`
		DatabaseDir string `json:"database_dir"`
	}
	ImportAliases struct {
		DataIA     string
		DatabaseIA string
		ProjectIA  string
	}
)

func New(file string) (s Settings, e error) {
	src, e := ioutil.ReadFile(file)
	if e != nil {
		return s, e
	}

	if e = json.Unmarshal(src, &s); e != nil {
		return s, e
	}

	if e = s.aliasingImports(); e != nil {
		return s, e
	}

	s.SqlDriver = strings.ToLower(s.SqlDriver)
	return s, e
}

func (settings *Settings) aliasingImports() error {
	prefix := filepath.Join("go", "src")

	if settings.DataDir == "" || settings.ProjectDir == "" {
		return errors.New("Не установлены пути ")
	}

	if strings.Contains(settings.ProjectDir, prefix) {
		settings.ProjectIA = extractImport(settings.ProjectDir)
	}
	if strings.Contains(settings.DataDir, prefix) {
		settings.DataIA = extractImport(settings.DataDir)
	}
	if strings.Contains(settings.DatabaseDir, prefix) {
		settings.DatabaseIA = extractImport(settings.DatabaseDir)
	}

	return nil
}

func extractImport(path string) string {
	b, e := len(path)-6, len(path)

	for i := 0; i < len(path); i++ {
		if b < 0 || e < 0 {
			return ""
		}

		if path[b:e] == filepath.Join("go", "src") {
			return changeSeparator(path[e+1:])
		}

		b, e = b-1, e-1
	}

	return ""
}

func changeSeparator(path string) string {
	var newpath strings.Builder
	pathways := strings.Split(path, string(os.PathSeparator))
	for k, pw := range pathways {
		newpath.WriteString(pw)
		if k < len(pathways)-1 {
			newpath.WriteString("/")
		}
	}
	return newpath.String()
}

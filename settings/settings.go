package settings

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

const (
	MySQL      = "MySQL"
	PostgreSQL = "PostgreSQL"
)

type Settings struct {
	Path `json:"path"`
	ImportAliases
	SqlDriver  string `json:"sql_driver"`
	AutoDelete bool   `json:"auto_delete"`
}

type Path struct {
	ProjectDir string `json:"project_dir"`
	DataDir    string `json:"data_dir"`
	StorageDir string `json:"storage_dir"`
}

type ImportAliases struct {
	DataIA    string
	StorageIA string
	ProjectIA string
}

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

	return s, e
}

func (settings *Settings) aliasingImports() error {
	if settings.DataDir == "" || settings.ProjectDir == "" {
		return errors.New("Не установлены пути ")
	}

	if strings.Contains(settings.ProjectDir, "go/src") {
		settings.ProjectIA = extractImport(settings.ProjectDir)
	}
	if strings.Contains(settings.DataDir, "go/src") {
		settings.DataIA = extractImport(settings.DataDir)
	}
	if strings.Contains(settings.StorageDir, "go/src") {
		settings.StorageIA = extractImport(settings.StorageDir)
	}

	return nil
}

func extractImport(path string) string {
	b, e := len(path)-6, len(path)

	for i := 0; i < len(path); i++ {
		if b < 0 || e < 0 {
			return ""
		}

		if path[b:e] == "go/src" {
			return path[e+1:]
		}

		b, e = b-1, e-1
	}

	return ""
}

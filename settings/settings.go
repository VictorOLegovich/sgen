package settings

import (
	"errors"
	"strings"
)

const (
	MySQL      = "MySQL"
	PostgreSQL = "PostgreSQL"
)

type Settings struct {
	Path
	ImportAliases
	SqlDriver   string
	PackageMode int
}

type Path struct {
	ProjectDir string
	DataDir    string
	StorageDir string
}

type ImportAliases struct {
	DataImportAlias    string
	StorageImportAlias string
	ProjectImportAlias string
}

func (settings *Settings) AliasingImports() error {
	if settings.DataDir == "" || settings.ProjectDir == "" {
		return errors.New("Не установлены пути ")
	}

	if strings.Contains(settings.ProjectDir, "go/src") {
		settings.ProjectImportAlias = extractImport(settings.ProjectDir)
	}
	if strings.Contains(settings.DataDir, "go/src") {
		settings.DataImportAlias = extractImport(settings.DataDir)
	}
	if strings.Contains(settings.StorageDir, "go/src") {
		settings.StorageImportAlias = extractImport(settings.StorageDir)
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
			return path[e:]
		}

		b, e = b-1, e-1
	}

	return ""
}

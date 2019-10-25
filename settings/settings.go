package settings

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
	ImportAlias        string
}

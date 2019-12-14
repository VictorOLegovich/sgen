package query_builder

type QbFile struct {
	Name, Src string
}

func Files() []QbFile {
	return []QbFile{
		{"delete.go", DELETE},
		{"insert.go", INSERT},
		{"select.go", SELECT},
		{"update.go", UPDATE},
		{"utils.go", UTILS},
		{"query_builder.go", QueryBuilder},
	}
}

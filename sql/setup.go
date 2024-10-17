package sql

import (
	"embed"
	"io"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Migrations() []string {
	migrationsDir := "migrations"
	migrationsFiles, err := migrationsFS.ReadDir(migrationsDir)
	if err != nil {
		return []string{}
	}
	slices.SortFunc(migrationsFiles, func(a, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	migrations := make([]string, len(migrationsFiles))
	for i, file := range migrationsFiles {
		fn := filepath.Join(migrationsDir, file.Name())
		f, err := migrationsFS.Open(fn)
		if err != nil {
			return []string{}
		}
		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			return []string{}
		}

		migrations[i] = string(content)
	}
	return migrations
}

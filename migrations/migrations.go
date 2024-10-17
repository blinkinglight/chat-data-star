package migrations

import (
	"database/sql"
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/johejo/golang-migrate-extra/source/iofs"
	_ "zombiezen.com/go/sqlite"
)

//go:embed sql/*
var FS embed.FS

func Migrate(file string) {
	d, err := iofs.New(FS, "sql")
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite", file)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "db", driver)

	if err != nil {
		log.Fatal(err)
	}
	_ = m

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Print(err)
	}

}

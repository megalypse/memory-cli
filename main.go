package main

import (
	"os"

	"github.com/megalypse/memory_cli/cmd/cli"
	db2 "github.com/megalypse/memory_cli/factory/db"
	"github.com/megalypse/memory_cli/internal/memorygroup"
	"github.com/pressly/goose/v3"
)

func main() {
	err := startup()
	if err != nil {
		os.Exit(1)
	}
}

func startup() error {
	env := "PROD"
	if val := os.Getenv("MaRY_ENV"); val != "" {
		env = val
	}

	db, err := db2.GetDB()
	if err != nil {
		return err
	}

	_ = memorygroup.GetRepositorySqlLite(db)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	goose.SetBaseFS(migrationsFS)

	if env == "DEV" {
		err = goose.Reset(db, "migrations")
		if err != nil {
			return err
		}
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		return err
	}

	return cli.Execute()
}

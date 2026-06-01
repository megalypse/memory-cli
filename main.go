package main

import (
	"os"

	"github.com/megalypse/memory_cli/cmd/cli"
	db2 "github.com/megalypse/memory_cli/factory/db"
	"github.com/pressly/goose/v3"
)

func main() {
	err := startup()
	if err != nil {
		os.Exit(1)
	}
}

func startup() error {
	db, err := db2.GetDB()
	if err != nil {
		return err
	}

	goose.SetBaseFS(migrationsFS)
	err = goose.Up(db, "migrations")
	if err != nil {
		return err
	}

	return cli.Execute()
}

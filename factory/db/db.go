package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

var GetDB = sync.OnceValues(getDB)

func getDB() (*sql.DB, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(dir, "memory_cli")

	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}
	dbPath := filepath.Join(appDir, "memory.db")
	return sql.Open("sqlite", dbPath)
}

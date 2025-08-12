package database

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"path/filepath"
	"sync"
	"watercolormc/internal"
	"watercolormc/internal/utils"
)

type Database struct {
	Client *sql.DB
}

var (
	dbInstance *Database
	once       sync.Once
)

// Init initializes the singleton Database instance once
func Init() error {
	if dbInstance != nil {
		zap.L().Warn("database already initialized")
		return nil
	}

	dataDir := utils.ExpandHome(internal.WatercolorDataDirectory)

	createErr := utils.CreateIfNotExists(dataDir)
	if createErr != nil {
		return createErr
	}

	dataSourceName := filepath.Join(dataDir, internal.DatabaseName)

	var err error
	once.Do(func() {
		db, e := sql.Open("sqlite3", dataSourceName)
		if e != nil {
			err = e
			return
		}
		if e = db.Ping(); e != nil {
			err = e
			return
		}
		dbInstance = &Database{Client: db}
	})

	return err
}

func SetupSchema() error {
	if dbInstance == nil {
		return errors.New("database not initialized")
	}

	schema := `
	CREATE TABLE IF NOT EXISTS servers (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		port INTEGER NOT NULL,
		host TEXT NOT NULL,
	    description TEXT DEFAULT '',
	    version TEXT NOT NULL,
	    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := dbInstance.Client.Exec(schema)
	return err
}

// Get returns the singleton Database instance
func Get() *Database {
	return dbInstance
}

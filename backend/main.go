package backend

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/mszalewicz/frosk/trace"
)

type Backend struct {
	DB     *sql.DB
	logger *slog.Logger
}

func Initialize(applicationDB string, logger *slog.Logger) (*Backend, error) {

	db, err := sql.Open("sqlite3", applicationDB)
	if err != nil {
		logger.Error(fmt.Errorf("Could not initializa db file: %w", err).Error(), slog.String("trace", trace.CurrentWindow()))
		return nil, err
	}

	backend := Backend{DB: db, logger: logger}
	return &backend, nil
}

func (backend *Backend) CreateStructure() error {

	const create_passwords_table = `
		CREATE TABLE IF NOT EXISTS passwords (
	       id INTEGER PRIMARY KEY AUTOINCREMENT,
	       name TEXT UNIQUE NOT NULL,
	       password TEXT NOT NULL,
	       created_at TEXT NULL,
	       updated_at TEXT NULL
	   ) STRICT;
	`

	_, err := backend.DB.Exec(create_passwords_table)

	if err != nil {
		backend.logger.Error(fmt.Errorf("Error during creating passwords table: %w", err).Error(), slog.String("trace", trace.CurrentWindow()))
		return err
	}

	const create_master_table = `
		CREATE TABLE IF NOT EXISTS master (
	    	id INTEGER PRIMARY KEY AUTOINCREMENT,
		    password TEXT UNIQUE NOT NULL,
		    secret_key TEXT UNIQUE NOT NULL,
		    iv TEXT UNIQUE NOT NULL,
		    created_at TEXT NULL,
		    updated_at TEXT NULL
		) STRICT;
	`

	_, err = backend.DB.Exec(create_master_table)

	if err != nil {
		backend.logger.Error(fmt.Errorf("Error during creating passwords table: %w", err).Error(), slog.String("trace", trace.CurrentWindow()))
		return err
	}

	return nil
}

func (backend *Backend) CountMasterEntries() (int, error) {
	var numberOfEntriesInMaster int
	row := backend.DB.QueryRow("SELECT COUNT(*) FROM master")

	err := row.Scan(numberOfEntriesInMaster)

	if err != nil {
		backend.logger.Error(fmt.Errorf("Query counting number of entries in master table: %w", err).Error(), slog.String("trace", trace.CurrentWindow()))
		return 0, err
	}

	return numberOfEntriesInMaster, nil
}

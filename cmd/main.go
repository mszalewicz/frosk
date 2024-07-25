package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/mszalewicz/frosk/backend"
	"github.com/mszalewicz/frosk/trace"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// TODO: set db path such that it will much OS native path scheme:
	// 		Mac:       ~/Library/Applications Support/frosk/log
	// 		Windows:   C:\Users\<username>\AppData\Local\frosk\log
	// 		Linux:     /var/lib/frosk/log
	const logPath string = "./cmd/log"
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		// TODO implement GUI response
		log.Fatal(fmt.Errorf("Error during opening log file: %w", err))
	}
	defer logFile.Close()
	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	// TODO: set db path such that it will much OS native path scheme:
	// 		Mac:       ~/Library/Applications Support/frosk/application.sqlite
	// 		Windows:   C:\Users\<username>\AppData\Local\frosk\application.sqlite
	// 		Linux:     /var/lib/frosk/application.sqlite
	const applicationDB string = "./sql/application.sqlite"

	backend, err := backend.Initialize(applicationDB, logger)

	if err != nil {
		logger.Error(fmt.Errorf("Error during initilization of database: %w", err).Error(), slog.String("trace", trace.CurrentWindow()))
		// TODO implement GUI response
	}

	err = backend.CreateStructure()

	if err != nil {
		logger.Error(fmt.Errorf("Error during creation of db schema: %w", err).Error(), slog.String("trace", trace.CurrentWindow()))
		// TODO implement GUI response
	}

	// n, err := backend.CountMasterEntries()

}

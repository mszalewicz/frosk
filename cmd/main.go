package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/mszalewicz/frosk/backend"
	"github.com/mszalewicz/frosk/helpers"

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

	loggerArgs := &slog.HandlerOptions{AddSource: true}
	logger := slog.New(slog.NewJSONHandler(logFile, loggerArgs))
	slog.SetDefault(logger)

	// TODO: set db path such that it will much OS native path scheme:
	// 		Mac:       ~/Library/Applications Support/frosk/application.sqlite
	// 		Windows:   C:\Users\<username>\AppData\Local\frosk\application.sqlite
	// 		Linux:     /var/lib/frosk/application.sqlite
	const applicationDB string = "./sql/application.sqlite"

	// errToHandleInGUI is an error that was logged in the function that returns it
	// It is returned to indicate that something went wrong and should be reflected in GUI state
	backend, errToHandleInGUI := backend.Initialize(applicationDB)

	if errToHandleInGUI != nil {
		// TODO implement GUI response
	}

	errToHandleInGUI = backend.CreateStructure()

	if errToHandleInGUI != nil {
		// TODO implement GUI response
	}

	// n, err := backend.CountMasterEntries()
	fmt.Println(helpers.RandString(150, true))
	os.Exit(1)
	// backend.InsertMaster("test")o

	// ""

	//    '2025-05-29 14:16:00'

}

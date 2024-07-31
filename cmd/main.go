package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"

	server "github.com/mszalewicz/frosk/backend"
	"golang.org/x/crypto/bcrypt"

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

	localDevLog := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// TODO: set db path such that it will much OS native path scheme:
	// 		Mac:       ~/Library/Applications Support/frosk/application.sqlite
	// 		Windows:   C:\Users\<username>\AppData\Local\frosk\application.sqlite
	// 		Linux:     /var/lib/frosk/application.sqlite
	const applicationDB string = "./sql/application.sqlite"

	// errToHandleInGUI is an error that was logged in the function that returns it
	// It is returned to indicate that something went wrong and should be reflected in GUI state
	backend, errToHandleInGUI := server.Initialize(applicationDB)

	if errToHandleInGUI != nil {
		fmt.Println(errToHandleInGUI)
		// TODO implement GUI response
	}

	localDevLog.Debug("Creating database from schema...")
	errToHandleInGUI = backend.CreateStructure()

	if errToHandleInGUI != nil {
		fmt.Println(errToHandleInGUI)
		// TODO implement GUI response
	}

	localDevLog.Debug("Checking number of entries in master table...")
	numberOfEntriesInMasterTable, err := backend.CountMasterEntries()

	if numberOfEntriesInMasterTable == 0 {
		// TODO: get master password from GUI
		// TODO: if master password < 1, show appriopriate GUI message
		localDevLog.Debug("Initializing master table entry...")
		errToHandleInGUI := backend.InitMaster("placeholder")
		if errToHandleInGUI != nil {
			fmt.Println(errToHandleInGUI)
			// TODO: implement GUI response
		}
	}

	localDevLog.Debug("Creating password entry...")
	errToHandleInGUI = backend.EncryptPasswordEntry("service name t", "password t", "username t", "placeholder")

	if errToHandleInGUI != nil {
		switch {
		case errors.Is(errToHandleInGUI, server.EmptyPassword):
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		case errors.Is(errToHandleInGUI, server.EmptyServiceName):
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		case errors.Is(errToHandleInGUI, server.EmptyMasterPassord):
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		case errors.Is(errToHandleInGUI, server.EmptyUsername):
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		case errors.Is(errToHandleInGUI, bcrypt.ErrMismatchedHashAndPassword): // Check for authentication
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		}
	}
}

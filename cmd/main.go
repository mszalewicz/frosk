package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	server "github.com/mszalewicz/frosk/backend"
	"github.com/mszalewicz/frosk/gui"
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

	localDevLog := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	localDevLog.Debug("Creating database from schema...")
	errToHandleInGUI = backend.CreateStructure()

	if errToHandleInGUI != nil {
		fmt.Println(errToHandleInGUI)
		// TODO implement GUI response
	}

	// GUI development ---------------------------------------

	go func() {
		window := new(app.Window)
		window.Option(app.Title("frosk"))
		window.Option(app.Size(unit.Dp(450), unit.Dp(800)))
		window.Option(app.MinSize(unit.Dp(350), unit.Dp(350)))
		window.Option(app.Decorated(false))

		err := gui.HandleMainWindow(window, backend)

		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()

	// TODO: Move below logic to window handlers: ---------------------------------------------------------

	serviceName := "google"
	pass := "supersecretpass"
	username := "bestestuser"
	masterpass := "placeholder"

	errToHandleInGUI = backend.EncryptPasswordEntry("test", "test", "test", masterpass) // TODO: get info from GUI

	localDevLog.Debug("Creating password entry...")
	errToHandleInGUI = backend.EncryptPasswordEntry(serviceName, pass, username, masterpass) // TODO: get info from GUI

	if errToHandleInGUI != nil {
		switch {
		case errors.Is(errToHandleInGUI, server.EmptyPassword):
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		case errors.Is(errToHandleInGUI, server.EmptyUsername):
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		case errors.Is(errToHandleInGUI, server.EmptyServiceName):
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		case errors.Is(errToHandleInGUI, server.EmptyMasterPassord):
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		case errors.Is(errToHandleInGUI, server.ServiceNameAlreadyTaken):
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		case errors.Is(errToHandleInGUI, bcrypt.ErrMismatchedHashAndPassword): // Check for authentication
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		default:
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		}
	}

	localDevLog.Debug("Getting password entry...")
	passwordEntry, errToHandleInGUI := backend.DecryptPasswordEntry(serviceName, masterpass)

	if errToHandleInGUI != nil {
		switch {
		case errors.Is(errToHandleInGUI, sql.ErrNoRows):
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		case errors.Is(errToHandleInGUI, bcrypt.ErrMismatchedHashAndPassword): // Check for authentication
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		default:
			// TODO: implement GUI response
			localDevLog.Debug(errToHandleInGUI.Error())
		}
	}

	fmt.Printf("Password entry - service name: %s | username: %s | password: %s\n", passwordEntry.ServiceName, passwordEntry.Username, passwordEntry.Password)

}

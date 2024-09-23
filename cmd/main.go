package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"gioui.org/app"
	"gioui.org/unit"
	server "github.com/mszalewicz/frosk/backend"
	"github.com/mszalewicz/frosk/gui"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// TODO: set db path such that it will much OS native path scheme:
	// 		Mac:       ~/Library/Applications Support/frosk/log
	// 		Windows:   C:\Users\<username>\AppData\Local\frosk\log
	// 		Linux:     /var/lib/frosk/log
	current_os := runtime.GOOS
	logPath := ""

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	switch current_os {
	case "darwin":
		appDirectory := filepath.Join(usr.HomeDir, "/Library/Application Support/frosk/")
		logPath = filepath.Join(appDirectory, "log")

		_, err = os.Stat(appDirectory)
		if os.IsNotExist(err) {
			err := os.MkdirAll(appDirectory, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
		} else if err != nil {
			fmt.Println("Error checking directory:", err)
			return
		}
	case "windows":
		appDirectory := filepath.Join(usr.HomeDir, "AppData\\Local\\frosk")
		logPath = filepath.Join(appDirectory, "log")

		_, err = os.Stat(appDirectory)
		if os.IsNotExist(err) {
			err := os.MkdirAll(appDirectory, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
		} else if err != nil {
			fmt.Println("Error checking directory:", err)
			return
		}
	}

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

	// TODO: check if db already created
	localDevLog.Debug("Create database from schema if not present...")
	errToHandleInGUI = backend.CreateStructure()

	if errToHandleInGUI != nil {
		fmt.Println(errToHandleInGUI)
		// TODO implement GUI response
	}

	// Start GUI
	go func() {
		window := new(app.Window)
		window.Option(app.Title("frosk"))
		window.Option(app.Size(unit.Dp(450), unit.Dp(800)))
		window.Option(app.MinSize(unit.Dp(350), unit.Dp(350)))
		window.Option(app.Decorated(false))

		err := gui.HandleMainWindow(window, backend)

		if err != nil {
			slog.Error(err.Error())
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

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
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	server "github.com/mszalewicz/frosk/backend"
	"github.com/mszalewicz/frosk/gui"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//      App directory scheme under different OS:
	// 		Mac:       ~/Library/Applications Support/frosk/log
	// 		Windows:   C:\Users\<username>\AppData\Local\frosk\log
	// 		Linux:     /var/lib/frosk/log

	current_os := runtime.GOOS
	logPath := ""
	applicationDBPath := ""
	appDirectory := ""

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	switch current_os {
	case "darwin":
		appDirectory = filepath.Join(usr.HomeDir, "/Library/Application Support/frosk/")
		logPath = filepath.Join(appDirectory, "log")
		applicationDBPath = filepath.Join(appDirectory, "application.sqlite")

	case "windows":
		appDirectory := filepath.Join(usr.HomeDir, "AppData\\Local\\frosk")
		logPath = filepath.Join(appDirectory, "log")
		applicationDBPath = filepath.Join(appDirectory, "application.sqlite")

	case "linux":
		appDirectory := "/var/lib/frosk"
		logPath = filepath.Join(appDirectory, "log")
		applicationDBPath = filepath.Join(appDirectory, "application.sqlite")
	}

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

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		slog.Error("Could not create log file.", "error", err)
		var ops op.Ops
		theme := material.NewTheme()

		go func() {
			window := new(app.Window)
			window.Option(app.Title("frosk"))
			window.Option(app.Size(unit.Dp(450), unit.Dp(800)))
			window.Option(app.MinSize(unit.Dp(350), unit.Dp(350)))
			window.Option(app.Decorated(false))
			gui.ErrorWindow(&ops, window, theme, "Application could not create log file. Please report bug. Details can be found in the logs.")
		}()

		app.Main()
	}
	defer logFile.Close()

	loggerArgs := &slog.HandlerOptions{AddSource: true}
	logger := slog.New(slog.NewJSONHandler(logFile, loggerArgs))
	slog.SetDefault(logger)

	backend, errToHandleInGUI := server.Initialize(applicationDBPath)

	if errToHandleInGUI != nil {
		slog.Error("Could not create local db.", "error", errToHandleInGUI)
		var ops op.Ops
		theme := material.NewTheme()

		go func() {
			window := new(app.Window)
			window.Option(app.Title("frosk"))
			window.Option(app.Size(unit.Dp(450), unit.Dp(800)))
			window.Option(app.MinSize(unit.Dp(350), unit.Dp(350)))
			window.Option(app.Decorated(false))
			gui.ErrorWindow(&ops, window, theme, "Application could not create local database. Please report bug. Details can be found in the logs.")
		}()

		app.Main()
	}

	localDevLog := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	localDevLog.Debug("Create database from schema if not present...")
	errToHandleInGUI = backend.CreateStructure()

	if errToHandleInGUI != nil {
		slog.Error("Could not bootstrap db from schema.", "error", errToHandleInGUI)
		var ops op.Ops
		theme := material.NewTheme()

		go func() {
			window := new(app.Window)
			window.Option(app.Title("frosk"))
			window.Option(app.Size(unit.Dp(450), unit.Dp(800)))
			window.Option(app.MinSize(unit.Dp(350), unit.Dp(350)))
			window.Option(app.Decorated(false))
			gui.ErrorWindow(&ops, window, theme, "Application could not create database structure. Please report bug. Details can be found in the logs.")
		}()

		app.Main()
	}

	// Start GUI
	localDevLog.Debug("Starting GUI...")
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

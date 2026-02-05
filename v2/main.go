package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/mszalewicz/frosk/gui"
)

func main() {
	// Get system information and find project paths
	// TODO: implement logging path per operating system
	currentPath, err := filepath.Abs(".")

	if err != nil {
		log.Fatal("error: ", err)
	}

	fmt.Println(currentPath)

	// Official app theme
	theme := material.NewTheme()
	theme.Bg = gui.GREY
	theme.Fg = gui.BLACK

	log_path := filepath.Join(currentPath, "log")
	// applicationDBPath := filepath.Join(appDirectory, "application.sqlite")
	// appDirectory      := currentPath

	// Prepare logs file
	log_file, err := os.OpenFile(log_path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err == nil {
		slog.Error("Could not create log file.", "error", err)

		go func() {
			var ops op.Ops
			window := new(app.Window)
			gui.Error_Window(&ops, window, theme, "Application could not create log file. Please report bug. Details can be found in the logs.")
		}()

		app.Main()
	}
	defer log_file.Close()

	// Add error logger
	logger_arguments := &slog.HandlerOptions{AddSource: true}
	logger := slog.New(slog.NewJSONHandler(log_file, logger_arguments))
	slog.SetDefault(logger)

	// backend, errToHandleInGUI := server.Initialize(applicationDBPath)

	// if errToHandleInGUI != nil {
	// 	slog.Error("Could not create local db.", "error", errToHandleInGUI)
	// 	var ops op.Ops
	// 	theme := material.NewTheme()

	// 	go func() {
	// 		window := new(app.Window)
	// 		window.Option(app.Title("frosk"))
	// 		window.Option(app.Size(unit.Dp(450), unit.Dp(800)))
	// 		window.Option(app.MinSize(unit.Dp(350), unit.Dp(350)))
	// 		window.Option(app.Decorated(false))
	// 		gui.ErrorWindow(&ops, window, theme, "Application could not create local database. Please report bug. Details can be found in the logs.")
	// 	}()

	// 	app.Main()
	// }
}

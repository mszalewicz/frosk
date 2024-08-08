package gui

import (
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"os"
	"time"

	server "github.com/mszalewicz/frosk/backend"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// var green = color.NRGBA{R: 52, G: 235, B: 131, A: 255}
var green = color.NRGBA{R: 100, G: 196, B: 166, A: 255}
var orange = color.NRGBA{R: 235, G: 178, B: 10, A: 100}
var red = color.NRGBA{R: 238, G: 78, B: 78, A: 255}
var purple = color.NRGBA{R: 161, G: 100, B: 196, A: 255}

// 108, 235, 169 / 201/

// Creates list entry components
func createPasswordEntryListLineComponents(serviceName string, theme *material.Theme) []layout.FlexChild {
	const buttonSize = 12

	var openBtnWidget widget.Clickable
	openBtn := material.Button(theme, &openBtnWidget, "OPEN")
	openBtn.Color = color.NRGBA{R: 0, B: 0, G: 0, A: 255}
	openBtn.Background = green
	// openBtn.Background = color.NRGBA{R: 67, G: 168, B: 84, A: 255}
	openBtn.TextSize = unit.Sp(buttonSize)
	openBtn.Font.Weight = font.Bold

	var deleteBtnWidget widget.Clickable
	deleteBtn := material.Button(theme, &deleteBtnWidget, "DELETE")
	deleteBtn.Color = color.NRGBA{R: 0, B: 0, G: 0, A: 255}
	deleteBtn.Background = purple
	// deleteBtn.Background = color.NRGBA{R: 235, G: 64, B: 52, A: 255}
	deleteBtn.TextSize = unit.Sp(buttonSize)
	deleteBtn.Font.Weight = font.Bold

	var btnMargin = layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(0)}
	var labelMargin = layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(0)}

	serviceFlexChild := layout.Flexed(
		1,
		func(gtx layout.Context) layout.Dimensions {
			return labelMargin.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {

					serviceNameLabel := material.Label(theme, unit.Sp(25), serviceName)
					// serviceNameLabel.Font.Weight = font.Bold
					serviceNameLabel.MaxLines = 1

					return serviceNameLabel.Layout(gtx)
				},
			)
		},
	)

	openBtnFlexChild := layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			return btnMargin.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {
					return openBtn.Layout(gtx)
				},
			)
		},
	)

	deleteBtnFlexChild := layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			return btnMargin.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {
					return deleteBtn.Layout(gtx)
				},
			)
		},
	)

	return []layout.FlexChild{serviceFlexChild, openBtnFlexChild, deleteBtnFlexChild}
}

// Creastes and populates GUI list container from password entries components
func constructPasswordEntriesList(passwordEntries [][]layout.FlexChild, passwordEntriesList *layout.List, margin layout.Inset) layout.FlexChild {
	return layout.Flexed(
		1,
		func(gtx layout.Context) layout.Dimensions {
			return margin.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {
					return passwordEntriesList.Layout(
						gtx,
						len(passwordEntries),
						func(gtx layout.Context, i int) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical}.Layout(
								gtx,
								layout.Rigid(
									func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, passwordEntries[i]...)
									},
								),
								horizontalDivider(),
							)
						},
					)
				},
			)
		},
	)
}

// Init horizontal line divider
func horizontalDivider() layout.FlexChild {
	return layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			height := unit.Dp(1)
			line := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Dp(height))
			paint.FillShape(gtx.Ops, color.NRGBA{A: 40}, clip.Rect(line).Op())
			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Dp(height))}
		},
	)
}

// func loadingScreen()

func HandleMainWindow(window *app.Window, backend *server.Backend) error {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	theme.Bg = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
	// theme.Bg = color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	initialRender := true
	// masterPasswordSet := false
	margin := layout.Inset{Top: unit.Dp(15), Bottom: unit.Dp(15), Left: unit.Dp(15), Right: unit.Dp(15)}

	var ops op.Ops
	var newPasswordEntryWidget widget.Clickable

	testServices := []string{"super long service name label", "test of language support: część", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank"}

	passwordEntriesList := &layout.List{Axis: layout.Vertical}
	passwordEntries := [][]layout.FlexChild{}

	for _, serviceName := range testServices {
		passwordEntries = append(passwordEntries, createPasswordEntryListLineComponents(serviceName, theme))
	}

	localDevLog := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	localDevLog.Debug("Checking number of entries in master table...")
	numberOfEntriesInMasterTable, errToHandleInGUI := backend.CountMasterEntries()

	if errToHandleInGUI != nil {
		localDevLog.Debug(errToHandleInGUI.Error())
	}

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

	// loader := material.Loader(theme)
	// loader.

	const loadWindowSize = 300
	window.Option(app.Size(unit.Dp(loadWindowSize*3), unit.Dp(loadWindowSize)))
	window.Option(app.MaxSize(unit.Dp(loadWindowSize*3), unit.Dp(loadWindowSize)))
	window.Option(app.MinSize(unit.Dp(loadWindowSize*3), unit.Dp(loadWindowSize)))

	// loaderMargin := layout.Inset{
	// 	Left:   unit.Dp(30),
	// 	Right:  unit.Dp(30),
	// 	Bottom: unit.Dp(30),
	// 	Top:    unit.Dp(30),
	// }

	centerWindows := true

	// loadingProgressChan := make(chan float32)
	// progress := float32(0)

	// go func() {
	// 	for {
	// 		time.Sleep(time.Second / 60)
	// 		loadingProgressChan <- 0.0055
	// 	}
	// }()

	// go func() {
	// 	for p := range loadingProgressChan {
	// 		if progress < 1 {
	// 			progress += p
	// 			window.Invalidate()
	// 		}
	// 	}
	// }()

	var text string
	var buttontest widget.Clickable
	textChannel := make(chan string)

	go func() {
		option := 0

		for {
			if option == 0 {
				textChannel <- ".  "
				option += 1
			} else if option == 1 {
				textChannel <- ".. "
				option += 1
			} else {
				textChannel <- "..."
				option = 0
			}

			time.Sleep(time.Second / 2)
			window.Invalidate()
		}
	}()

	startTime := time.Now()

InitLoop:
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
				// layout.Flexed(1,
				// 	func(gtx layout.Context) layout.Dimensions {
				// 		return loaderMargin.Layout(gtx,
				// 			func(gtx layout.Context) layout.Dimensions {
				// 				// progress += <-loadingProgressChan
				// 				progressCircle := material.ProgressCircle(theme, progress)
				// 				progressCircle.Color = color.NRGBA{R: 123, G: 78, B: 90, A: 255}
				// 				return progressCircle.Layout(gtx)
				// 			},
				// 		)
				// 	},
				// ),
				layout.Flexed(1,
					func(gtx layout.Context) layout.Dimensions {

						// test := rand.Intn(8)

						// if test <= 2 {
						// 	text = "Initializing..."
						// } else if 3 <= test && test <= 5 {
						// 	text = "Initializing.  "
						// } else {
						// 	text = "Initializing.. "
						// }

						select {
						case msg := <-textChannel:
							text = msg
						default:
						}

						btn := material.Button(theme, &buttontest, "First time setup in progress"+text)
						btn.TextSize = unit.Sp(40)
						btn.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
						btn.CornerRadius = unit.Dp(0)
						// btn.Background = color.NRGBA{R: 52, G: 235, B: 131, A: 255}
						btn.Background = orange
						// 52, 189, 235
						// 52, 235, 131

						return btn.Layout(gtx)
					},
				),
			)

			if centerWindows {
				window.Perform(system.ActionCenter)
				// window.Option(app.Size(loadWindowSize, loadWindowSize))
				// window.Invalidate()

				// const loadWindowSize = 300
				// window.Option(app.Size(unit.Dp(loadWindowSize), unit.Dp(loadWindowSize)))
				// window.Option(app.MaxSize(unit.Dp(loadWindowSize), unit.Dp(loadWindowSize)))
				// window.Option(app.MinSize(unit.Dp(loadWindowSize), unit.Dp(loadWindowSize)))
				centerWindows = false
			}

			if time.Since(startTime).Seconds() > 1 {
				// window.Option(app.Title("frosk"))
				window.Option(app.MinSize(unit.Dp(350), unit.Dp(350)))
				window.Option(app.MaxSize(unit.Dp(2000), unit.Dp(2000)))
				window.Option(app.Size(unit.Dp(450), unit.Dp(800)))
				window.Option(app.Decorated(true))
				// window.Perform(system.ActionCenter)
				// window.Invalidate()
				break InitLoop
			}

			e.Frame(gtx.Ops)

		}
	}

	// progressCircle := material.ProgressCircle(theme, 2)
	// progressCircle.Color = color.NRGBA{}
	// progressCircle.Layout(gtx)

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:

			// TODO; implement remembering last window size
			// fmt.Println("x: ", e.Size.X, " y: ", e.Size.Y, " conversion:", e.Metric.PxPerDp)

			gtx := app.NewContext(&ops, e)

			if newPasswordEntryWidget.Clicked(gtx) {
				fmt.Println("test")
			}

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(
				gtx,
				constructPasswordEntriesList(passwordEntries, passwordEntriesList, margin),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return margin.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								newPasswordEntry := material.Button(theme, &newPasswordEntryWidget, "NEW")
								newPasswordEntry.Background = color.NRGBA{R: 30, G: 30, B: 30, A: 255}
								newPasswordEntry.TextSize = unit.Sp(25)
								newPasswordEntry.Font.Weight = font.Bold
								return newPasswordEntry.Layout(gtx)
							},
						)
					},
				),
			)

			e.Frame(gtx.Ops)

			if initialRender {
				window.Perform(system.ActionCenter)
				initialRender = !initialRender
			}
		}
	}
}

package gui

import (
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"os"

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

type PasswordEntriesGUI struct {
	serviceName    string
	guiListElement []layout.FlexChild
}

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
func constructPasswordEntriesList(passwordEntries *[]PasswordEntriesGUI, passwordEntriesList *layout.List, margin layout.Inset) layout.FlexChild {
	return layout.Flexed(
		1,
		func(gtx layout.Context) layout.Dimensions {
			return margin.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {
					return passwordEntriesList.Layout(
						gtx,
						len(*passwordEntries),
						func(gtx layout.Context, i int) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical}.Layout(
								gtx,
								layout.Rigid(
									func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, (*passwordEntries)[i].guiListElement...)
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

func errorWindow(ops *op.Ops, window *app.Window, theme *material.Theme, errorMsg string) error {
	ResizeWindowInfo(window)
	centerWindow := true

	errConfirmWidget := new(widget.Clickable)
	errListContainer := &widget.List{List: layout.List{Axis: layout.Vertical, Alignment: layout.Start}}

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(ops, e)

			if errConfirmWidget.Clicked(gtx) {
				os.Exit(0)
			}

			InfoWindowWidget(&gtx, theme, errConfirmWidget, errListContainer, errorMsg)

			if centerWindow {
				window.Perform(system.ActionCenter)
				centerWindow = !centerWindow
			}

			e.Frame(gtx.Ops)
		}
	}
}

func HandleMainWindow(window *app.Window, backend *server.Backend) error {
	errChan := make(chan error)

	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	theme.Bg = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
	theme.ContrastBg = red
	// theme.Bg = color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	// masterPasswordSet := false
	margin := layout.Inset{Top: unit.Dp(15), Bottom: unit.Dp(15), Left: unit.Dp(15), Right: unit.Dp(15)}

	var ops op.Ops
	var newPasswordEntryWidget widget.Clickable

	// testServices := []string{"super long service name label", "test of language support: część", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank"}

	localDevLog := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	localDevLog.Debug("Checking number of entries in master table...")
	numberOfEntriesInMasterTable, errToHandleInGUI := backend.CountMasterEntries()

	if errToHandleInGUI != nil {
		localDevLog.Debug(errToHandleInGUI.Error())
	}

	centerWindow := true

	ResizeWindowInitialSetup(window)

	passwordInput := new(widget.Editor)
	passwordInput.SingleLine = true
	passwordInput.Mask = '*'
	passwordInput.Filter = alphabet

	passwordInputRepeat := new(widget.Editor)
	passwordInputRepeat.SingleLine = true
	passwordInputRepeat.Mask = '*'
	passwordInputRepeat.Filter = alphabet

	confirmBtnWidget := new(widget.Clickable)
	showHideWidget := new(widget.Clickable)

	initialSetup := InitialSetup{
		passwordInput:       passwordInput,
		passwordInputRepeat: passwordInputRepeat,
		confirmBtnWidget:    confirmBtnWidget,
		showHidWidget:       showHideWidget,
		borderColor:         black,
	}

	// Get master password info during firt use of application
	if numberOfEntriesInMasterTable == 0 {
	InitialSetupWindowMarker:
		for {
			switch e := window.Event().(type) {
			case app.DestroyEvent:
				return e.Err

			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				// Show/hide input
				if initialSetup.showHidWidget.Clicked(gtx) {
					switch {
					case initialSetup.passwordInput.Mask == rune(0):
						initialSetup.passwordInput.Mask = '*'
						initialSetup.passwordInputRepeat.Mask = '*'
					default:
						initialSetup.passwordInput.Mask = rune(0)
						initialSetup.passwordInputRepeat.Mask = rune(0)
					}
				}

				// Check if password are non empty and match
				if confirmBtnWidget.Clicked(gtx) {
					switch {
					case passwordInput.Len() > 0 && passwordInputRepeat.Len() > 0:
						if passwordInput.Text() == passwordInputRepeat.Text() {
							localDevLog.Debug("Checking number of entries in master table...")

							go func() {
								err := backend.InitMaster(passwordInput.Text())
								errChan <- err
								return
							}()

							break InitialSetupWindowMarker
						}
					default:
						initialSetup.borderColor = red
					}
				}

				// Handle master password input events
			CheckInputEventMarker:
				for {
					_, passwordInputEvent := initialSetup.passwordInput.Update(gtx)
					_, passwordInputRepeatEvent := initialSetup.passwordInputRepeat.Update(gtx)

					inputEventOccured := passwordInputEvent || passwordInputRepeatEvent

					switch {
					case inputEventOccured && initialSetup.passwordInput.Len() > 0 && initialSetup.passwordInputRepeat.Len() > 0 && initialSetup.passwordInput.Text() != initialSetup.passwordInputRepeat.Text():
						initialSetup.borderColor = red
					case inputEventOccured:
						initialSetup.borderColor = black
					default:
						break CheckInputEventMarker
					}
				}

				InitialSetupWidget(&gtx, theme, &initialSetup)

				if centerWindow {
					window.Perform(system.ActionCenter)
					centerWindow = !centerWindow
				}

				e.Frame(gtx.Ops)
			}
		}

		ResizeWindowLoad(window)
		centerWindow = true
		errOccured := false

		// Show loader during master password save action
	InitialLoadMarker:
		for {
			switch e := window.Event().(type) {
			case app.DestroyEvent:
				return e.Err

			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				select {
				case err := <-errChan:
					if err != nil {
						errWrapped := fmt.Errorf("Could not save master password in database: %w", err)
						slog.Error(errWrapped.Error())
						errOccured = true
					}
					break InitialLoadMarker
				default:
					// paint.Fill(gtx.Ops, black)
					LoadWidget(&gtx, theme)
				}

				if centerWindow {
					window.Perform(system.ActionCenter)
					centerWindow = !centerWindow
				}

				e.Frame(gtx.Ops)
			}
		}

		// Handle master password save error
		if errOccured {
			errorWindow(&ops, window, theme, "Error during saving master password.")
		}
	}

	centerWindow = true

	for {
		services, err := backend.GetPasswordEntriesList()

		if err != nil {
			errorWindow(&ops, window, theme, "Could not load password entries.")
		}

		passwordEntriesList := &layout.List{Axis: layout.Vertical}
		passwordEntries := []PasswordEntriesGUI{}
		fmt.Println(services)

		for _, serviceName := range services {
			listElement := createPasswordEntryListLineComponents(serviceName, theme)
			passwordEntries = append(passwordEntries, PasswordEntriesGUI{serviceName: serviceName, guiListElement: listElement})
		}

		ResizeWindowPasswordEntriesList(window)

	ShowListMarker:
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
					break ShowListMarker
				}

				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceStart,
				}.Layout(
					gtx,
					constructPasswordEntriesList(&passwordEntries, passwordEntriesList, margin),
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

				if centerWindow {
					window.Perform(system.ActionCenter)
					centerWindow = !centerWindow
				}
			}
		}

		// Entry new password
		for {
			switch e := window.Event().(type) {
			case app.DestroyEvent:
				return e.Err

			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				// TODO: Get input for new password entry

				e.Frame(gtx.Ops)
			}
		}

	}
}

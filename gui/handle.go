package gui

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"os"
	"strconv"
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

	numberOfEntriesInMasterTable, errToHandleInGUI := backend.CountMasterEntries()

	if errToHandleInGUI != nil {
		errorWindow(&ops, window, theme, "Fatal error when running application. Please consult logs.")
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

	firstTimeShowingInitialSetupWindow := true

	// Get master password info during firt use of application
	if numberOfEntriesInMasterTable == 0 {
		if firstTimeShowingInitialSetupWindow {
			go func() {
				time.Sleep(time.Second / 20)
				window.Invalidate()
				return
			}()
			firstTimeShowingInitialSetupWindow = !firstTimeShowingInitialSetupWindow
		}

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
	firstTimeShowing := true

	for {
		services, err := backend.GetPasswordEntriesList()

		if err != nil {
			errorWindow(&ops, window, theme, "Could not load password entries.")
		}

		passwordEntriesList := &layout.List{Axis: layout.Vertical}
		passwordEntries := []PasswordEntriesGUI{}

		for _, serviceName := range services {
			listElement := createPasswordEntryListLineComponents(serviceName, theme)
			passwordEntries = append(passwordEntries, PasswordEntriesGUI{serviceName: serviceName, guiListElement: listElement})
		}

		// Schedule invalidate in seperate gorotuine to redraw window after initial show.
		// For some reason gio do not paint correct layout / elements sizes on the first show after centering.
		if firstTimeShowing {
			go func() {
				time.Sleep(time.Second / 20)
				window.Invalidate()
				return
			}()
			firstTimeShowing = !firstTimeShowing
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

				if centerWindow {
					window.Perform(system.ActionCenter)
					centerWindow = !centerWindow
				}

				e.Frame(gtx.Ops)
			}
		}

		err = InputNewPassword(window, &ops, backend, theme)

		if err != nil {
			errorWindow(&ops, window, theme, "Error occured during password saving. Please check logs.")
		}

		centerWindow = true
	}
}

type Information struct {
	text  string
	color color.NRGBA
}

func InputNewPassword(window *app.Window, ops *op.Ops, backend *server.Backend, theme *material.Theme) error {
	var centerWindow bool = true
	var inserted bool = true

	ResizeWindowNewPasswordInsert(window)

	masterPassword := new(widget.Editor)
	masterPassword.SingleLine = true
	masterPassword.Mask = '*'
	masterPassword.Filter = alphabet

	serviceName := new(widget.Editor)
	serviceName.SingleLine = true
	serviceName.Mask = '*'
	serviceName.Filter = alphabet

	username := new(widget.Editor)
	username.SingleLine = true
	username.Mask = '*'
	username.Filter = alphabet

	password := new(widget.Editor)
	password.SingleLine = true
	password.Mask = '*'
	password.Filter = alphabet

	confirmBtnWidget := new(widget.Clickable)
	showHideWidget := new(widget.Clickable)

	newPasswordView := NewPasswordView{
		masterPassword:   masterPassword,
		serviceName:      serviceName,
		username:         username,
		password:         password,
		confirmBtnWidget: confirmBtnWidget,
		showHidWidget:    showHideWidget,
		borderColor:      black,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	passwordLength := ""
	countLetterChan := make(chan int)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			countLetterChan <- newPasswordView.password.Len()
			time.Sleep(time.Second / 8)
		}
	}(ctx)

	info := Information{"Provide Master Password to authenticate. Fill out form to save password for a service.", purple}
	tryingToInsertPassword := false

	type InsertPasswordEntryOperation struct {
		error     error
		didInsert bool
		msg       string
	}

	insertPasswordOperationChan := make(chan InsertPasswordEntryOperation)

	// Draw
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(ops, e)

			select {
			case insertOperation := <-insertPasswordOperationChan:
				if insertOperation.error != nil {
					switch err := insertOperation.error; {
					case errors.Is(err, server.ServiceNameAlreadyTaken), errors.Is(err, server.MasterPasswordDoNotMatch):
						info.text = insertOperation.msg
						info.color = red
						tryingToInsertPassword = false
						ResizeWindowNewPasswordInsert(window)
						window.Perform(system.ActionCenter)
					default:
						return err
					}
				}
				if insertOperation.didInsert {
					return nil
				}
			default:
			}

		CheckConfirmButtonClickMarker:
			if confirmBtnWidget.Clicked(gtx) {
				info.text = ""
				inputProblem := false

				if len(newPasswordView.masterPassword.Text()) == 0 {
					info.text += "Master Password is empty. "
					info.color = red
					inputProblem = true
				}
				if len(newPasswordView.username.Text()) == 0 {
					info.text += "Username is empty. "
					info.color = red
					inputProblem = true
				}
				if len(newPasswordView.serviceName.Text()) == 0 {
					info.text += "Service name is empty. "
					info.color = red
					inputProblem = true
				}
				if len(newPasswordView.password.Text()) == 0 {
					info.text += "Password is empty. "
					info.color = red
					inputProblem = true
				}

				if inputProblem {
					goto CheckConfirmButtonClickMarker
				}

				go func() {
					masterPasswordMatch, err := backend.CmpMasterPassword(newPasswordView.masterPassword.Text())

					if !masterPasswordMatch && err == nil {
						insertPasswordOperationChan <- InsertPasswordEntryOperation{server.MasterPasswordDoNotMatch, !inserted, "Master Password is incorrect."}
						return
					}

					if err != nil {
						fmt.Println("wut")
					}

					err = backend.EncryptPasswordEntry(
						newPasswordView.serviceName.Text(),
						newPasswordView.password.Text(),
						newPasswordView.username.Text(),
						newPasswordView.masterPassword.Text(),
					)

					if err != nil {
						if errors.Is(err, server.ServiceNameAlreadyTaken) {
							insertPasswordOperationChan <- InsertPasswordEntryOperation{err, !inserted, "Service name is already taken. Choose another name."}
						} else {
							insertPasswordOperationChan <- InsertPasswordEntryOperation{err, !inserted, "Unspecified error occured. Check error description."}
						}
						return
					}

					insertPasswordOperationChan <- InsertPasswordEntryOperation{nil, inserted, ""}
				}()

				tryingToInsertPassword = true
				ResizeWindowLoad(window)
				window.Perform(system.ActionCenter)
			}

			if showHideWidget.Clicked(gtx) {
				switch {
				case newPasswordView.masterPassword.Mask == rune(0):
					newPasswordView.masterPassword.Mask = '*'
					newPasswordView.serviceName.Mask = '*'
					newPasswordView.username.Mask = '*'
					newPasswordView.password.Mask = '*'
				default:
					newPasswordView.masterPassword.Mask = rune(0)
					newPasswordView.serviceName.Mask = rune(0)
					newPasswordView.username.Mask = rune(0)
					newPasswordView.password.Mask = rune(0)
				}
			}

			select {
			case count := <-countLetterChan:
				passwordLength = strconv.Itoa(count)
			default:
			}

			window.Invalidate()

			if tryingToInsertPassword {
				LoadWidget(&gtx, theme)
			} else {
				InsertNewPasswordWidget(&gtx, theme, &newPasswordView, passwordLength, info)
			}

			if centerWindow {
				window.Perform(system.ActionCenter)
				centerWindow = !centerWindow
			}

			e.Frame(gtx.Ops)
		}
	}
}

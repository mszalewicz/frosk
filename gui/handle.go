package gui

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"sync"

	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"time"

	server "github.com/mszalewicz/frosk/backend"
	"github.com/mszalewicz/frosk/helpers"

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
	"github.com/sahilm/fuzzy"
)

type PasswordEntriesGUI struct {
	serviceName     string
	guiListElement  []layout.FlexChild
	openBtnWidget   *widget.Clickable
	deleteBtnWidget *widget.Clickable
}

type DecryptionPackage struct {
	err           error
	passwordEntry server.PasswordEntry
}

// Creates list entry components
func createPasswordEntryListLineComponents(serviceName string, theme *material.Theme) ([]layout.FlexChild, *widget.Clickable, *widget.Clickable) {
	const buttonSize = 12

	var openBtnWidget widget.Clickable
	openBtn := material.Button(theme, &openBtnWidget, "OPEN")
	openBtn.Color = black
	openBtn.Background = grey_light
	openBtn.TextSize = unit.Sp(buttonSize)
	openBtn.Font.Weight = font.Medium
	openBtn.Font.Typeface = "Verdana, monospace"
	// openBtn.Font.Style

	var deleteBtnWidget widget.Clickable
	deleteBtn := material.Button(theme, &deleteBtnWidget, "DELETE")
	deleteBtn.Color = black
	deleteBtn.Background = grey_light
	deleteBtn.TextSize = unit.Sp(buttonSize)
	deleteBtn.Font.Weight = font.Medium
	deleteBtn.Font.Typeface = "Verdana, monospace"

	var btnMargin = layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(0)}
	var labelMargin = layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(0)}

	serviceFlexChild := layout.Flexed(
		1,
		func(gtx layout.Context) layout.Dimensions {
			return labelMargin.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {

					serviceNameLabel := material.Label(theme, unit.Sp(25), serviceName)
					serviceNameLabel.Font.Weight = font.Normal
					serviceNameLabel.Font.Typeface = "Verdana, monospace"
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
					border := widget.Border{Color: charcoal, CornerRadius: unit.Dp(4), Width: unit.Dp(0)}
					return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Dp(0)).Layout(gtx, openBtn.Layout)
					})
				},
			)
		},
	)

	deleteBtnFlexChild := layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			return btnMargin.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {
					border := widget.Border{Color: charcoal, CornerRadius: unit.Dp(4), Width: unit.Dp(0)}
					return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Dp(0)).Layout(gtx, deleteBtn.Layout)
					})
				},
			)
		},
	)

	return []layout.FlexChild{serviceFlexChild, openBtnFlexChild, deleteBtnFlexChild}, &openBtnWidget, &deleteBtnWidget
}

// Creastes and populates GUI list container from password entries components
func constructPasswordEntriesList(passwordEntries *[]PasswordEntriesGUI, passwordEntriesList *layout.List, margin layout.Inset, mutex *sync.Mutex) layout.FlexChild {

	mutex.Lock()
	localEntries := *passwordEntries
	mutex.Unlock()

	return layout.Flexed(
		1,
		func(gtx layout.Context) layout.Dimensions {
			return margin.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {
					return passwordEntriesList.Layout(
						gtx,
						len(localEntries),
						func(gtx layout.Context, i int) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical}.Layout(
								gtx,
								layout.Rigid(
									func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, (localEntries)[i].guiListElement...)
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

func emptyDivider() layout.FlexChild {
	return layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			height := unit.Dp(20)
			line := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Dp(height))
			paint.FillShape(gtx.Ops, white, clip.Rect(line).Op())
			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Dp(height))}
		},
	)
}

func ErrorWindow(ops *op.Ops, window *app.Window, theme *material.Theme, errorMsg string) error {
	ResizeWindowInfo(window)
	centerWindow := true

	errConfirmWidget := new(widget.Clickable)
	errListContainer := &widget.List{List: layout.List{Axis: layout.Vertical, Alignment: layout.Start}}

	go func() {
		for range 3 {
			time.Sleep(time.Second / 20)
			window.Invalidate()
		}
		return
	}()

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(1)
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(ops, e)

			if errConfirmWidget.Clicked(gtx) {
				os.Exit(1)
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
	theme.Bg = grey_light
	theme.ContrastBg = grey_light

	margin := layout.Inset{Top: unit.Dp(15), Bottom: unit.Dp(15), Left: unit.Dp(15), Right: unit.Dp(15)}

	var ops op.Ops
	var newPasswordEntryWidget widget.Clickable

	numberOfEntriesInMasterTable, errToHandleInGUI := backend.CountMasterEntries()

	if errToHandleInGUI != nil {
		ErrorWindow(&ops, window, theme, "Fatal error when running application. Please consult logs.")
	}

	centerWindow := true

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

	var searchInput widget.Editor
	searchInput.SingleLine = true

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
		ResizeWindowInitialSetup(window)

		if firstTimeShowingInitialSetupWindow {
			go func() {
				for range 3 {
					time.Sleep(time.Second / 20)
					window.Invalidate()
				}
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
			var errorWindowOps op.Ops
			ErrorWindow(&errorWindowOps, window, theme, "Error during saving master password.")
		}
	}

	refreshChan := make(chan bool, 1)
	centerWindow = true
	var passwordListOps op.Ops
	var mutex sync.Mutex

PasswordViewMarker:
	for {
		services, err := backend.GetPasswordEntriesList()

		if err != nil {
			var errorWindowOps op.Ops
			ErrorWindow(&errorWindowOps, window, theme, "Could not load password entries.")
		}

		passwordEntriesList := &layout.List{Axis: layout.Vertical}
		passwordEntries := make([]PasswordEntriesGUI, 0, len(services))

		for _, serviceName := range services {
			listElement, openBtnWidget, deleteBtnWidget := createPasswordEntryListLineComponents(serviceName, theme)
			passwordEntries = append(passwordEntries, PasswordEntriesGUI{serviceName: serviceName, guiListElement: listElement, openBtnWidget: openBtnWidget, deleteBtnWidget: deleteBtnWidget})
		}

		fullSetOfPasswordEntries := passwordEntries

		var serviceNames []string

		for _, entry := range passwordEntries {
			serviceNames = append(serviceNames, entry.serviceName)
		}

		// Schedule invalidate in seperate gorotuine to redraw window after initial show after resizing + centering.
		// For some reason gio do not paint correct layout / elements sizes on the first show after resizing + centering.
		go func() {
			for range 3 {
				time.Sleep(time.Second / 20)
				window.Invalidate()
			}
			return
		}()

		ResizeWindowPasswordEntriesList(window)

		// ShowListMarker:
		for {
			switch e := window.Event().(type) {
			case app.DestroyEvent:
				return e.Err

			case app.FrameEvent:
				gtx := app.NewContext(&passwordListOps, e)

				select {
				case shouldRefresh := <-refreshChan:
					_ = shouldRefresh
					searchInput.SetText("")
					goto PasswordViewMarker
				default:
				}

				{ // Fuzzy search over service names list
					event, ok := searchInput.Update(gtx)
					if ok {
						if _, ok := event.(widget.ChangeEvent); ok {
							if searchInput.Text() != "" {
								go func() {
									query := searchInput.Text()
									matches := fuzzy.Find(query, serviceNames)

									var validEntries []PasswordEntriesGUI
									for _, match := range matches {
										for _, entry := range fullSetOfPasswordEntries {
											if entry.serviceName == match.Str {
												validEntries = append(validEntries, entry)
											}
										}
									}

									mutex.Lock()
									passwordEntries = validEntries
									mutex.Unlock()
								}()
							} else {
								go func() {
									mutex.Lock()
									passwordEntries = fullSetOfPasswordEntries
									mutex.Unlock()
								}()
							}

						}
					}
				}

				for _, passwordEntryInfo := range passwordEntries {
					if passwordEntryInfo.openBtnWidget.Clicked(gtx) {
						go authenticateAndShowPassword(backend, theme, passwordEntryInfo.serviceName)
					}

					if passwordEntryInfo.deleteBtnWidget.Clicked(gtx) {
						go confirmDeletion(backend, theme, passwordEntryInfo.serviceName, refreshChan)
					}

				}

				if newPasswordEntryWidget.Clicked(gtx) {
					go func() {
						var newPasswordEntryOps op.Ops
						newPasswordWindow := new(app.Window)
						ResizeWindowNewPasswordInsert(newPasswordWindow)
						err = InputNewPassword(newPasswordWindow, &newPasswordEntryOps, backend, theme, refreshChan)

						if err != nil {
							var errorWindowOps op.Ops
							ErrorWindow(&errorWindowOps, newPasswordWindow, theme, "Error occured during password saving. Please check logs.")
						}
					}()
				}

				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceStart,
				}.Layout(
					gtx,

					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return DrawSearchInput(gtx, theme, &searchInput, 130)
					}),

					constructPasswordEntriesList(&passwordEntries, passwordEntriesList, margin, &mutex),
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							return margin.Layout(gtx,
								func(gtx layout.Context) layout.Dimensions {
									newPasswordEntry := material.Button(theme, &newPasswordEntryWidget, "NEW")
									newPasswordEntry.Background = charcoal
									newPasswordEntry.TextSize = unit.Sp(25)
									newPasswordEntry.Font.Weight = font.SemiBold
									newPasswordEntry.Font.Typeface = "Verdana, monospace"

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
	}
}

func authenticateAndShowPassword(backend *server.Backend, theme *material.Theme, serviceName string) {
	var (
		centerWindow                  bool = true
		alreadyDecrypted              bool = false
		authenticate                  widget.Clickable
		cancel                        widget.Clickable
		showHideUsername              widget.Clickable
		showHidePassword              widget.Clickable
		masterPasswordGUI             widget.Editor
		usernameGUI                   widget.Editor
		passwordGUI                   widget.Editor
		textCheckMsg                  string
		passwordEditorBackgroundColor color.NRGBA = grey
	)

	masterPasswordGUI.SingleLine = true
	masterPasswordGUI.Mask = '*'

	usernameGUI.ReadOnly = true
	usernameGUI.SingleLine = true
	usernameGUI.Mask = '*'

	passwordGUI.ReadOnly = true
	passwordGUI.SingleLine = true
	passwordGUI.Mask = '*'

	ops := new(op.Ops)
	window := new(app.Window)
	ResizeDecryptionWindow(window)

	confirmDecryptionChan := make(chan DecryptionPackage, 2)
	closeLoaderChan := make(chan bool, 2)

	// Schedule invalidate in seperate gorotuine to redraw window after initial show after resizing + centering.
	// For some reason gio do not paint correct layout / elements sizes on the first show after resizing + centering.
	go func() {
		for range 3 {
			time.Sleep(time.Second / 20)
			window.Invalidate()
		}
		return
	}()

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return
		case app.FrameEvent:
			gtx := app.NewContext(ops, e)
			select {
			case decryptPackage := <-confirmDecryptionChan:
				closeLoaderChan <- true
				switch decryptErr := decryptPackage.err; {
				case decryptErr == nil:
					passwordEditorBackgroundColor = white

					usernameGUI.ReadOnly = false
					if usernameGUI.Text() != decryptPackage.passwordEntry.Password {
						usernameGUI.SetText(decryptPackage.passwordEntry.Username)
					}

					passwordGUI.ReadOnly = false
					if passwordGUI.Text() != decryptPackage.passwordEntry.Password {
						passwordGUI.SetText(decryptPackage.passwordEntry.Password)
					}

					alreadyDecrypted = !alreadyDecrypted
				case errors.Is(decryptErr, server.MasterPasswordDoNotMatch):
					textCheckMsg = " - incorrect password."
					passwordEditorBackgroundColor = red
				default:
				}
			default:
			}

			if authenticate.Clicked(gtx) && !alreadyDecrypted {
				if len(masterPasswordGUI.Text()) == 0 {
					textCheckMsg = " - empty, please enter password"
				} else {
					opsLoading := new(op.Ops)
					loaderWindow := new(app.Window)
					ResizeWindowLoad(loaderWindow)

					textCheckMsg = ""

					masterPassword := masterPasswordGUI.Text()
					go tryPasswordDecryption(backend, window, confirmDecryptionChan, &serviceName, &masterPassword)
					go showLoading(opsLoading, loaderWindow, theme, closeLoaderChan)
				}
			}

			if cancel.Clicked(gtx) {
				window.Perform(system.ActionClose)
			}

			if showHidePassword.Clicked(gtx) {
				if passwordGUI.ReadOnly != true {
					switch {
					case passwordGUI.Mask == rune(0):
						passwordGUI.Mask = '*'
					default:
						passwordGUI.Mask = rune(0)
					}
				}
			}

			if showHideUsername.Clicked(gtx) {
				if usernameGUI.ReadOnly != true {
					switch {
					case usernameGUI.Mask == rune(0):
						usernameGUI.Mask = '*'
					default:
						usernameGUI.Mask = rune(0)
					}
				}
			}

			ManagePasswordDecryptionWidget(&gtx, theme, &serviceName, &textCheckMsg, &authenticate, &cancel, &showHideUsername, &showHidePassword, &masterPasswordGUI, &usernameGUI, &passwordGUI, &passwordEditorBackgroundColor)

			if centerWindow {
				window.Perform(system.ActionCenter)
				centerWindow = !centerWindow
			}

			e.Frame(gtx.Ops)
		}
	}
}

func tryPasswordDecryption(backend *server.Backend, window *app.Window, confirmDecryptionChan chan DecryptionPackage, serviceName *string, masterPassword *string) {
	_, err := backend.CmpMasterPassword(*masterPassword)

	if err != nil {
		confirmDecryptionChan <- DecryptionPackage{err: err, passwordEntry: server.PasswordEntry{}}
		window.Invalidate()
		return
	}

	passwordEntry, err := backend.DecryptPasswordEntry(*serviceName, *masterPassword)

	if err != nil {
		confirmDecryptionChan <- DecryptionPackage{err: err, passwordEntry: server.PasswordEntry{}}
	} else {
		confirmDecryptionChan <- DecryptionPackage{err: nil, passwordEntry: passwordEntry}
	}

	window.Invalidate()
	return
}

func showLoading(ops *op.Ops, window *app.Window, theme *material.Theme, closeLoadingChan chan bool) {
	var centerWindow bool = true

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return

		case app.FrameEvent:
			gtx := app.NewContext(ops, e)

			select {
			case <-closeLoadingChan:
				window.Perform(system.ActionClose)
			default:
				LoadWidget(&gtx, theme)
			}

			if centerWindow {
				window.Perform(system.ActionCenter)
				centerWindow = !centerWindow
			}

			e.Frame(gtx.Ops)
		}
	}
}

func confirmDeletion(backend *server.Backend, theme *material.Theme, serviceName string, refreshChan chan bool) {
	var (
		deletePasswordEntry bool = true
		centerWindow        bool = true
		confirm             widget.Clickable
		deny                widget.Clickable
		maxW                unit.Dp = 650
		maxH                unit.Dp = 350
	)

	ops := new(op.Ops)
	window := new(app.Window)
	window.Option(app.Decorated(false))
	window.Option(app.MinSize(unit.Dp(maxW), unit.Dp(maxH)))
	window.Option(app.MaxSize(unit.Dp(maxW), unit.Dp(maxH)))
	window.Option(app.Size(unit.Dp(maxW), unit.Dp(maxH)))
	window.Option(app.Title("frosk"))

	go func() {
		time.Sleep(time.Second / 20)
		window.Invalidate()
		return
	}()

	{ // Window loop
		for {
			switch e := window.Event().(type) {
			case app.DestroyEvent:
				// return e.Err
				return

			case app.FrameEvent:
				gtx := app.NewContext(ops, e)

				{ // Choice whether to delete password or not
					if confirm.Clicked(gtx) {
						err := backend.DeletePasswordEntry(serviceName)
						if err != nil {
							errWrapped := fmt.Errorf("Error during deletion of password: %w", err)
							slog.Error(errWrapped.Error())
						}
						refreshChan <- deletePasswordEntry
						window.Perform(system.ActionClose)
					}

					if deny.Clicked(gtx) {
						window.Perform(system.ActionClose)
					}
				}

				{ // Paint
					paint.Fill(ops, black)
					ConfirmPasswordDeletionWidget(&gtx, theme, serviceName, &confirm, &deny)

					if centerWindow {
						window.Perform(system.ActionCenter)
						centerWindow = !centerWindow
					}

					e.Frame(gtx.Ops)
				}

			}
		}
	}

}

type Information struct {
	text  string
	color color.NRGBA
}

func InputNewPassword(window *app.Window, ops *op.Ops, backend *server.Backend, theme *material.Theme, refreshChan chan bool) error {
	var centerWindow bool = true
	var inserted bool = true

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
	smallRandWidget := new(widget.Clickable)
	mediumRandWidget := new(widget.Clickable)
	bigRandWidget := new(widget.Clickable)
	specialCharsWidget := new(widget.Clickable)
	specialCharsText := "Special Chars: ON"
	specialCharsColor := orange

	newPasswordView := NewPasswordView{
		masterPassword:           masterPassword,
		serviceName:              serviceName,
		username:                 username,
		password:                 password,
		confirmBtnWidget:         confirmBtnWidget,
		showHidWidget:            showHideWidget,
		smallRandWidget:          smallRandWidget,
		mediumRandWidget:         mediumRandWidget,
		bigRandWidget:            bigRandWidget,
		specialCharsSwitchWidget: specialCharsWidget,
		specialCharsSwitchText:   specialCharsText,
		specialCharsSwitchColor:  specialCharsColor,
		specialCharsFlag:         true,
		borderColor:              black,
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

	info := Information{"Provide Master Password to authenticate. Fill out form to save credentials for a service.", purple}
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
					refreshChan <- true
					window.Perform(system.ActionClose)
				}
			default:
			}

			if smallRandWidget.Clicked(gtx) {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				min := 20
				max := 25
				len := r.Intn(max-min+1) + min
				randomString := helpers.RandString(len, newPasswordView.specialCharsFlag)
				newPasswordView.password.SetText(randomString)
			}
			if mediumRandWidget.Clicked(gtx) {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				min := 25
				max := 40
				len := r.Intn(max-min+1) + min
				randomString := helpers.RandString(len, newPasswordView.specialCharsFlag)
				newPasswordView.password.SetText(randomString)
			}
			if bigRandWidget.Clicked(gtx) {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				min := 40
				max := 100
				len := r.Intn(max-min+1) + min
				randomString := helpers.RandString(len, newPasswordView.specialCharsFlag)
				newPasswordView.password.SetText(randomString)
			}

			if specialCharsWidget.Clicked(gtx) {
				if newPasswordView.specialCharsFlag == true {
					newPasswordView.specialCharsFlag = !newPasswordView.specialCharsFlag
					newPasswordView.specialCharsSwitchText = "Special Chars: OFF"
					newPasswordView.specialCharsSwitchColor = grey
				} else {
					newPasswordView.specialCharsFlag = !newPasswordView.specialCharsFlag
					newPasswordView.specialCharsSwitchText = "Special Chars: ON"
					newPasswordView.specialCharsSwitchColor = orange
				}
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

					if !masterPasswordMatch && errors.Is(err, server.MasterPasswordDoNotMatch) {
						insertPasswordOperationChan <- InsertPasswordEntryOperation{server.MasterPasswordDoNotMatch, !inserted, "Master Password is incorrect."}
						return
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

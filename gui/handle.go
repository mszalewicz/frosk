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
					paint.Fill(gtx.Ops, black)
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

	// var text string
	// var buttontest widget.Clickable
	// textChannel := make(chan string)

	// go func() {
	// 	option := 0

	// 	for {
	// 		if option == 0 {
	// 			textChannel <- ".  "
	// 			option += 1
	// 		} else if option == 1 {
	// 			textChannel <- ".. "
	// 			option += 1
	// 		} else {
	// 			textChannel <- "..."
	// 			option = 0
	// 		}

	// 		time.Sleep(time.Second / 2)
	// 		window.Invalidate()
	// 	}
	// }()

	// startTime := time.Now()

	// lorem := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aenean sit amet mi in diam gravida tincidunt. Sed sem elit, rhoncus lobortis diam sed, efficitur aliquam sem. Sed volutpat tempor nisi, et vehicula nisi vestibulum ac. Aliquam faucibus augue vel gravida viverra. Fusce tempus vel massa eu vehicula. Interdum et malesuada fames ac ante ipsum primis in faucibus. Quisque volutpat augue at nunc sodales viverra. Fusce metus justo, pharetra at odio vel, ornare venenatis dolor. Ut eget ipsum vel velit convallis pulvinar. Quisque sit amet nulla vitae ante suscipit sodales sit amet eget magna. Phasellus ornare at libero eget elementum. Morbi commodo ante sem, vitae suscipit libero molestie sit amet. Duis ac metus tincidunt nisi pretium pretium et sed dui. Curabitur et nunc velit. Aliquam interdum mi ligula, at dapibus sem efficitur vel.

	// 	In sagittis vel felis sed luctus. Ut vel gravida mauris, vel venenatis mauris. Maecenas ac vulputate elit. Maecenas eget lectus in ante sollicitudin placerat. Curabitur tempus pharetra lacinia. Nunc sem tellus, dapibus et ligula nec, lacinia euismod nulla. Nulla facilisi. Quisque varius gravida tortor ac tincidunt. Nulla facilisi. Donec blandit felis neque, ut rutrum lacus aliquam ac. Praesent dictum mattis dolor eu porta. Sed ac ullamcorper neque. Proin suscipit pretium libero vel auctor. Mauris eget ultricies nisi. Ut eget dolor at dui pharetra efficitur.

	// 	Sed consequat elit non mauris vehicula, et faucibus lectus vehicula. Maecenas aliquam laoreet nulla, nec mattis mauris tempus at. Sed nec leo vel nisi faucibus eleifend sed vitae mauris. Vivamus in gravida ante. Integer eget odio sollicitudin, porttitor risus non, facilisis dolor. Cras vel ex sagittis, viverra mauris id, vestibulum dolor. Suspendisse hendrerit augue sed leo sollicitudin pellentesque. Quisque suscipit augue at gravida accumsan. Pellentesque molestie pulvinar lorem eu maximus. Aliquam faucibus tellus sed odio tincidunt efficitur. Maecenas vitae nibh ac nulla condimentum consectetur. Proin porta sagittis elit et luctus. Sed hendrerit enim at ullamcorper aliquet.

	// 	Cras et suscipit leo, vitae accumsan sapien. Aliquam at elit lobortis, interdum ligula in, posuere metus. Nam et tristique ligula. Aliquam maximus turpis non sem malesuada, in pharetra urna viverra. Aenean sit amet magna erat. Sed dictum justo dolor, in mattis eros pretium a. Sed cursus orci mauris, vel efficitur nunc congue vitae. Etiam vehicula feugiat elit, vitae porttitor justo iaculis ac. In sit amet leo lectus. Morbi gravida purus enim, porttitor vulputate turpis tincidunt aliquam. Quisque et ex molestie, tempus augue quis, tempus massa.

	// 	Nam vitae egestas orci. Cras rutrum, velit at condimentum blandit, eros neque laoreet dui, vel fermentum velit mauris non velit. Praesent non nisl vel metus scelerisque posuere. Vivamus sapien eros, pharetra eget blandit dictum, mollis vel ligula. Vestibulum scelerisque lobortis tristique. Curabitur imperdiet porttitor condimentum. Interdum et malesuada fames ac ante ipsum primis in faucibus. Aenean quis suscipit lacus. Sed vestibulum congue tempus. Mauris facilisis, libero ac euismod eleifend, erat metus lobortis ex, non pulvinar felis tellus et orci. Etiam semper fringilla malesuada. Ut lectus lorem, tincidunt vel scelerisque vel, malesuada ut lacus. Integer eget nibh quam. Integer aliquam vitae arcu feugiat venenatis. Sed convallis porta nisl, et porta lorem vehicula ut.

	// 	Ut egestas mauris dapibus, faucibus mi ac, posuere nulla. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nunc tristique mauris sed leo aliquam cursus. Morbi vitae mi quis leo eleifend tristique ultrices ultrices urna. Integer dignissim ex metus, ac facilisis lacus maximus eu. Nunc imperdiet elit dapibus nisl efficitur ornare. Suspendisse mollis, purus et pulvinar faucibus, massa tortor mollis est, ut semper purus mauris nec ipsum. Sed varius dolor eu ultricies ultricies. Nam tincidunt, turpis id volutpat aliquet, ante dui suscipit mi, in ornare purus magna sed ligula. Nunc vitae quam sed ligula ullamcorper consectetur in sit amet tortor. Duis vitae risus justo. Donec commodo in eros at tincidunt. Sed vel lorem mattis, viverra nulla a, consectetur nisi. Morbi a risus sollicitudin sem porttitor pharetra non id nibh. Mauris nec condimentum arcu. Quisque congue bibendum sem et consectetur.

	// 	Donec consectetur enim sed felis dapibus consequat. Phasellus ultrices faucibus justo, ut sagittis eros luctus maximus. Integer eu mauris vel lorem suscipit luctus. Duis vel erat eget elit gravida lobortis. Maecenas enim ipsum, semper ut bibendum vitae, gravida vitae orci. Quisque ac orci lectus. Etiam ut viverra purus. Praesent ante massa, euismod eget ultrices eget, fringilla non risus. Aenean ultrices purus at arcu porta auctor. Donec non tellus risus. Nullam tincidunt bibendum erat, et euismod sapien pharetra ut.

	// 	Fusce commodo lacus tortor, non dignissim justo placerat ut. Nulla facilisi. Sed ut sapien ac orci suscipit pretium sit amet sit amet quam. Nam tincidunt felis id enim facilisis, sit amet lobortis erat ornare. Etiam et magna elementum, blandit diam quis, bibendum risus. Vestibulum rutrum a leo eu semper. Aenean sit amet venenatis purus, sed semper odio. Fusce sit amet enim at magna elementum facilisis. Quisque nec nisl ipsum. Duis mollis, neque at rhoncus rutrum, ante ante varius orci, ut iaculis felis lorem sed augue. Sed vehicula, felis ac blandit ultrices, leo odio posuere nulla, non tempor lorem orci eget sem. Etiam mollis ligula quis mi semper eleifend. In urna enim, lacinia sit amet odio et, aliquet pellentesque lacus.

	// 	Etiam pulvinar fermentum risus id commodo. Maecenas non neque turpis. Integer pretium ex vel ultrices gravida. Maecenas egestas fringilla ipsum, eu pharetra sapien sodales id. Vivamus cursus lectus sit amet sodales placerat. Mauris vitae porta arcu, vel condimentum metus. Donec placerat neque quis rhoncus mollis. Mauris pulvinar cursus justo at hendrerit. Phasellus odio nunc, eleifend ut dui quis, mollis volutpat mauris. Integer nunc metus, mattis eget eleifend ut, finibus nec sapien. Suspendisse potenti. Sed et ultrices est. Morbi semper orci sed risus volutpat rutrum. Nulla facilisi. Aenean maximus nisi eu mauris lobortis hendrerit.

	// 	Integer efficitur tincidunt massa, quis facilisis nibh tincidunt sed. Nam tincidunt nunc non fringilla venenatis. Nunc lobortis elit dolor, vitae posuere sem euismod ut. Ut vestibulum, urna id sagittis fermentum, justo massa porta diam, at dapibus lectus ipsum pretium odio. Proin pharetra maximus enim. Quisque at convallis risus, tempor tristique mauris. Pellentesque et scelerisque metus.

	// 	In hac habitasse platea dictumst. Sed quis sodales mauris, ut placerat ex. Cras nec auctor mi. Vivamus egestas erat sed bibendum ullamcorper. Sed tincidunt, massa eu iaculis accumsan, eros eros commodo felis, non dictum sapien nisi accumsan libero. Duis rutrum bibendum sapien nec tristique. Sed blandit condimentum neque finibus mattis. Nam imperdiet nunc erat, id laoreet nunc pellentesque vitae. Nam orci nisi, mattis id elementum a, scelerisque vitae ex. Pellentesque ullamcorper, felis ac efficitur scelerisque, odio elit vehicula eros, id commodo leo nulla sed nibh. Vivamus sapien justo, elementum ut elementum vel, scelerisque vel tellus. Curabitur ultrices scelerisque augue ac cursus. Fusce tempor sed lacus ut feugiat. Suspendisse auctor felis a velit malesuada tempor.

	// 	Nunc vel dapibus erat. Praesent velit justo, malesuada sit amet est ac, suscipit interdum nisi. Quisque quam ante, malesuada a lorem a, consequat porttitor justo. Ut ultrices cursus ipsum, quis tincidunt est blandit eget. Aliquam lacinia non ex et pretium. Donec tellus dui, placerat et ante ac, lacinia luctus nunc. Mauris porttitor, quam a venenatis accumsan, dui leo ultrices lectus, in hendrerit libero dolor non urna.

	// 	Aenean non molestie elit, id elementum urna. Integer fringilla nunc eu massa tempus volutpat vel in est. Suspendisse sed sollicitudin eros. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Fusce quam ante, porta et sapien cursus, feugiat aliquam sapien. Donec viverra, nisl at scelerisque bibendum, purus nulla molestie purus, eu consequat nisi quam ut nisl. Nunc neque est, accumsan a varius id, gravida cursus diam. Integer pretium ipsum nulla. Nam ullamcorper auctor tellus vitae vulputate. Vestibulum posuere tempor lobortis. Vestibulum a quam ut orci elementum sollicitudin. Cras varius non felis non eleifend. Aenean porttitor neque vitae dolor scelerisque blandit. Donec pulvinar bibendum orci at consectetur.

	// 	Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Donec interdum ornare gravida. Nulla purus augue, hendrerit in blandit convallis, laoreet at leo. Fusce sed imperdiet tortor, eu egestas nunc. In hac habitasse platea dictumst. Fusce ultricies sagittis nunc nec tristique. Duis blandit elit quis ipsum rutrum, lobortis feugiat nisl dignissim. Integer eget orci molestie, congue lacus ut, placerat metus. Proin pellentesque elit nisi, sit amet bibendum nulla blandit et. Phasellus venenatis mauris non lacus ullamcorper pellentesque. Vivamus in nunc non turpis elementum pellentesque.

	// 	Duis ullamcorper felis pulvinar lacus auctor, ut pretium nibh ornare. Maecenas vel luctus nunc, ultricies sollicitudin libero. Suspendisse vel malesuada justo, in cursus ligula. Suspendisse laoreet neque dui, vitae auctor ante lobortis vel. Cras tincidunt eget felis non malesuada. Sed viverra elit vitae sem tincidunt, vel elementum purus lacinia. Aliquam dictum consectetur justo non luctus. In hac habitasse platea dictumst.

	// 	Praesent porttitor tempus arcu. Mauris ex magna, hendrerit a elit a, tempor pulvinar leo. Donec tempor est a orci rutrum, vehicula congue turpis consectetur. Ut pretium tellus nibh, rutrum posuere mauris cursus vitae. Integer egestas consequat leo sit amet eleifend. Morbi dapibus dolor a enim pharetra, quis tristique metus aliquet. Suspendisse sollicitudin vel ligula nec elementum.

	// 	Vivamus et odio enim. Proin quis ultrices est, non vehicula justo. Nullam nec tincidunt augue. Aliquam aliquam ante non augue tempus aliquet ac et quam. Aenean gravida suscipit lacus ut ornare. In vestibulum dapibus sapien, in euismod urna. Duis luctus dignissim diam sed rutrum. Curabitur eget augue purus. Fusce suscipit nisi pellentesque, auctor lectus id, aliquam nisi.

	// 	Praesent rutrum libero at nulla pulvinar eleifend. Nulla feugiat gravida eros id rhoncus. Sed arcu libero, vehicula id lorem ac, finibus bibendum arcu. Donec posuere hendrerit rutrum. Vestibulum lacus nibh, pulvinar eget nulla ac, interdum volutpat quam. Sed eu efficitur sapien, a faucibus sapien. Curabitur nec faucibus lectus. Morbi suscipit gravida sem a mollis. Phasellus eget tortor magna. In id justo ac leo maximus vulputate a vel dui. Maecenas semper vitae elit nec egestas.
	// 	`

	// list := &widget.List{List: layout.List{Axis: layout.Vertical, Alignment: layout.Start}}
	// returnBtnWidget := new(widget.Clickable)
	// infoCentered := false

	// for {
	// 	switch e := window.Event().(type) {
	// 	case app.DestroyEvent:
	// 		return e.Err

	// 	case app.FrameEvent:
	// 		gtx := app.NewContext(&ops, e)
	// 		InfoWindowWidget(&gtx, theme, returnBtnWidget, list, &lorem)

	// 		if !infoCentered {
	// 			ResizeWindowInfo(window)
	// 			window.Perform(system.ActionCenter)
	// 			infoCentered = !infoCentered
	// 		}

	// 		e.Frame(gtx.Ops)
	// 	}
	// }

	// InitLoop:
	// 	for {
	// 		switch e := window.Event().(type) {
	// 		case app.DestroyEvent:
	// 			return e.Err

	// 		case app.FrameEvent:
	// 			gtx := app.NewContext(&ops, e)

	// 			layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
	// 				// layout.Flexed(1,
	// 				// 	func(gtx layout.Context) layout.Dimensions {
	// 				// 		return loaderMargin.Layout(gtx,
	// 				// 			func(gtx layout.Context) layout.Dimensions {
	// 				// 				// progress += <-loadingProgressChan
	// 				// 				progressCircle := material.ProgressCircle(theme, progress)
	// 				// 				progressCircle.Color = color.NRGBA{R: 123, G: 78, B: 90, A: 255}
	// 				// 				return progressCircle.Layout(gtx)
	// 				// 			},
	// 				// 		)
	// 				// 	},
	// 				// ),
	// 				layout.Flexed(1,
	// 					func(gtx layout.Context) layout.Dimensions {

	// 						// test := rand.Intn(8)

	// 						// if test <= 2 {
	// 						// 	text = "Initializing..."
	// 						// } else if 3 <= test && test <= 5 {
	// 						// 	text = "Initializing.  "
	// 						// } else {
	// 						// 	text = "Initializing.. "
	// 						// }

	// 						select {
	// 						case msg := <-textChannel:
	// 							text = msg
	// 						default:
	// 						}

	// 						btn := material.Button(theme, &buttontest, "First time setup in progress"+text)
	// 						btn.TextSize = unit.Sp(40)
	// 						btn.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	// 						btn.CornerRadius = unit.Dp(0)
	// 						// btn.Background = color.NRGBA{R: 52, G: 235, B: 131, A: 255}
	// 						btn.Background = orange
	// 						// 52, 189, 235
	// 						// 52, 235, 131

	// 						return btn.Layout(gtx)
	// 					},
	// 				),
	// 			)

	// 			if centerWindows {
	// 				window.Perform(system.ActionCenter)
	// 				// window.Option(app.Size(loadWindowSize, loadWindowSize))
	// 				// window.Invalidate()

	// 				// const loadWindowSize = 300
	// 				// window.Option(app.Size(unit.Dp(loadWindowSize), unit.Dp(loadWindowSize)))
	// 				// window.Option(app.MaxSize(unit.Dp(loadWindowSize), unit.Dp(loadWindowSize)))
	// 				// window.Option(app.MinSize(unit.Dp(loadWindowSize), unit.Dp(loadWindowSize)))
	// 				centerWindows = false
	// 			}

	// 			if time.Since(startTime).Seconds() > 1 {
	// 				// window.Option(app.Title("frosk"))
	// 				window.Option(app.MinSize(unit.Dp(350), unit.Dp(350)))
	// 				window.Option(app.MaxSize(unit.Dp(2000), unit.Dp(2000)))
	// 				window.Option(app.Size(unit.Dp(450), unit.Dp(800)))
	// 				window.Option(app.Decorated(true))
	// 				// window.Perform(system.ActionCenter)
	// 				// window.Invalidate()
	// 				break InitLoop
	// 			}

	// 			e.Frame(gtx.Ops)

	// 		}
	// 	}

	// progressCircle := material.ProgressCircle(theme, 2)
	// progressCircle.Color = color.NRGBA{}
	// progressCircle.Layout(gtx)

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

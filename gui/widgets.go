package gui

import (
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	beige        = color.NRGBA{R: 247, G: 239, B: 229, A: 255}
	black        = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	blue         = color.NRGBA{R: 150, G: 201, B: 244, A: 255}
	charcoal     = color.NRGBA{R: 55, G: 55, B: 55, A: 255}
	charcoal2    = color.NRGBA{R: 95, G: 95, B: 95, A: 255}
	green        = color.NRGBA{R: 100, G: 196, B: 166, A: 255}
	grey         = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	grey_light   = color.NRGBA{R: 235, G: 235, B: 235, A: 255}
	orange       = color.NRGBA{R: 235, G: 178, B: 10, A: 100}
	purple       = color.NRGBA{R: 161, G: 100, B: 196, A: 255}
	purple_light = color.NRGBA{R: 161, G: 100, B: 196, A: 120}
	red          = color.NRGBA{R: 238, G: 78, B: 78, A: 220}
	white        = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	textSize     = unit.Sp(30)
	appName      = "Vault"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789" + "!#$%&'~`(){}[]*+,-./:;<=>?@_|\"\\"

func ResizeWindowInfo(window *app.Window) {
	window.Option(app.Decorated(true))
	window.Option(app.MinSize(unit.Dp(300), unit.Dp(300)))
	window.Option(app.MaxSize(unit.Dp(450), unit.Dp(450)))
	window.Option(app.Size(unit.Dp(450), unit.Dp(450)))
	window.Option(app.Title(appName))
}

func InfoWindowWidget(gtx *layout.Context, theme *material.Theme, returnBtnWidget *widget.Clickable, list *widget.List, text string) {
	layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20), Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(
		*gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.End}.Layout(
				gtx,
				layout.Flexed(
					1,
					func(gtx layout.Context) layout.Dimensions {
						return list.Layout(
							gtx,
							1,
							func(gtx layout.Context, index int) layout.Dimensions {
								return material.Label(theme, unit.Sp(18), text).Layout(gtx)
							},
						)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						returnBtn := material.Button(theme, returnBtnWidget, "RETURN")
						returnBtn.Background = green
						returnBtn.Color = black
						returnBtn.Font.Weight = font.Bold
						return returnBtn.Layout(gtx)
					},
				),
			)
		},
	)
}

func ResizeWindowInitialSetup(window *app.Window) {
	window.Option(app.Decorated(true))
	window.Option(app.MinSize(unit.Dp(300), unit.Dp(300)))
	window.Option(app.MaxSize(unit.Dp(2000), unit.Dp(2000)))
	window.Option(app.Size(unit.Dp(850), unit.Dp(850)))
	window.Option(app.Title(appName))
}

type InitialSetup struct {
	passwordInput       *widget.Editor
	passwordInputRepeat *widget.Editor

	confirmBtnWidget *widget.Clickable
	showHidWidget    *widget.Clickable

	borderColor color.NRGBA
}

func InitialSetupWidget(gtx *layout.Context, theme *material.Theme, initialSetup *InitialSetup) {
	elementMargin := layout.Inset{Top: unit.Dp(17), Bottom: unit.Dp(17), Right: unit.Dp(10), Left: unit.Dp(10)}
	btnsMargin := layout.Inset{Top: unit.Dp(25), Bottom: unit.Dp(25), Right: unit.Dp(10), Left: unit.Dp(10)}

	layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(60), Right: unit.Dp(60)}.Layout(
		*gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceSides}.Layout(
				gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return elementMargin.Layout(
							gtx,
							func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceSides}.Layout(
									gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions { return material.H2(theme, "Initial Setup").Layout(gtx) }),
								)
							},
						)
					},
				),

				///TODO
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return elementMargin.Layout(
							gtx,
							func(gtx layout.Context) layout.Dimensions {
								label := material.Label(theme, unit.Sp(20), "Please enter your master password. A longer password provides stronger encryption.\n\nWarning: This password is the only key to your local vault. If lost, your data cannot be recovered and will be permanently inaccessible.")
								label.Color = blue
								label.Font.Weight = font.Medium
								label.Font.Typeface = "Verdana"
								return label.Layout(gtx)
							},
						)
					},
				),
				horizontalDivider(),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							return material.H6(theme, "Master Password:").Layout(gtx)
						},
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							inputMasterPassword := material.Editor(theme, initialSetup.passwordInput, "Enter password...")
							inputMasterPassword.TextSize = unit.Sp(20)
							inputMasterPassword.SelectionColor = blue
							// border := widget.Border{Color: initialSetup.borderColor, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}

							// return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(10)).Layout(gtx, inputMasterPassword.Layout)
							// })
						},
					)
				}),
				horizontalDivider(),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							return material.H6(theme, "Repeat Master Password:").Layout(gtx)
						},
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							inputMasterPasswordRepeat := material.Editor(theme, initialSetup.passwordInputRepeat, "Enter password...")
							inputMasterPasswordRepeat.TextSize = unit.Sp(20)
							inputMasterPasswordRepeat.SelectionColor = blue
							// border := widget.Border{Color: initialSetup.borderColor, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}

							// return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(10)).Layout(gtx, inputMasterPasswordRepeat.Layout)
							// })
						},
					)
				}),
				horizontalDivider(),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return btnsMargin.Layout(
							gtx,
							func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.End, Spacing: layout.SpaceStart}.Layout(
									gtx,
									layout.Rigid(
										func(gtx layout.Context) layout.Dimensions {
											return elementMargin.Layout(
												gtx,
												func(gtx layout.Context) layout.Dimensions {
													confirmBtn := material.Button(theme, initialSetup.confirmBtnWidget, "confirm")
													confirmBtn.Background = blue
													confirmBtn.Color = black
													confirmBtn.TextSize = unit.Sp(20)
													confirmBtn.Font.Weight = font.Medium
													confirmBtn.Font.Typeface = "Verdana"
													return confirmBtn.Layout(gtx)
												},
											)
										},
									),
									layout.Rigid(
										func(gtx layout.Context) layout.Dimensions {
											return elementMargin.Layout(
												gtx,
												func(gtx layout.Context) layout.Dimensions {
													confirmBtn := material.Button(theme, initialSetup.showHidWidget, "show/hide")
													confirmBtn.Background = grey_light
													confirmBtn.Color = black
													confirmBtn.TextSize = unit.Sp(20)
													confirmBtn.Font.Weight = font.Medium
													confirmBtn.Font.Typeface = "Verdana"
													return confirmBtn.Layout(gtx)
												},
											)
										},
									),
								)
							},
						)
					},
				),
			)
		},
	)
}

func ResizeWindowLoad(window *app.Window) {
	window.Option(app.Decorated(false))
	window.Option(app.MinSize(unit.Dp(300), unit.Dp(300)))
	window.Option(app.MaxSize(unit.Dp(300), unit.Dp(300)))
	window.Option(app.Size(unit.Dp(300), unit.Dp(300)))
}

func LoadWidget(gtx *layout.Context, theme *material.Theme) {
	layout.UniformInset(unit.Dp(20)).Layout(
		*gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(
				gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					loader := material.Loader(theme)
					loader.Color = black
					return loader.Layout(gtx)
				}),
			)
		},
	)
}

func ResizeWindowPasswordEntriesList(window *app.Window) {
	window.Option(app.Decorated(true))
	window.Option(app.MinSize(unit.Dp(300), unit.Dp(300)))
	// window.Option(app.MaxSize(unit.Dp(1_000), unit.Dp(1_000)))
	window.Option(app.Size(unit.Dp(500), unit.Dp(800)))
	window.Option(app.Title("Vault"))
}

func ResizeWindowNewPasswordInsert(window *app.Window) {
	window.Option(app.Decorated(true))
	window.Option(app.MinSize(unit.Dp(500), unit.Dp(800)))
	window.Option(app.MaxSize(unit.Dp(2000), unit.Dp(2000)))
	window.Option(app.Size(unit.Dp(750), unit.Dp(900)))
	window.Option(app.Title(appName))
}

type NewPasswordView struct {
	masterPassword *widget.Editor
	password       *widget.Editor
	serviceName    *widget.Editor
	username       *widget.Editor

	confirmBtnWidget         *widget.Clickable
	showHidWidget            *widget.Clickable
	smallRandWidget          *widget.Clickable
	mediumRandWidget         *widget.Clickable
	bigRandWidget            *widget.Clickable
	specialCharsSwitchWidget *widget.Clickable

	specialCharsSwitchText  string
	specialCharsSwitchColor color.NRGBA
	specialCharsFlag        bool

	borderColor color.NRGBA
}

func DrawSearchInput(gtx layout.Context, th *material.Theme, editor *widget.Editor, height int) layout.Dimensions {
	gtx.Constraints.Min.Y = height
	gtx.Constraints.Max.Y = height
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	margins := layout.Inset{
		Top:    unit.Dp(15),
		Bottom: unit.Dp(0),
		Left:   unit.Dp(20),
		Right:  unit.Dp(20),
	}

	return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widget.Border{
			Color:        color.NRGBA{A: 255},
			Width:        unit.Dp(1),
			CornerRadius: unit.Dp(4),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			bgColor := color.NRGBA{R: 245, G: 245, B: 245, A: 255}
			rr := gtx.Dp(unit.Dp(4)) // Match Border CornerRadius
			paint.FillShape(gtx.Ops, bgColor, clip.RRect{
				Rect: image.Rectangle{Max: gtx.Constraints.Min},
				SE:   rr, SW: rr, NW: rr, NE: rr,
			}.Op(gtx.Ops))

			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				editor := material.Editor(th, editor, "Search...")
				editor.TextSize = unit.Sp(28)
				editor.SelectionColor = purple_light
				editor.Font.Typeface = "Verdana, monospace"
				// Center text and hint horizontally
				return editor.Layout(gtx)
			})
		})
	})
}

func InsertNewPasswordWidget(gtx *layout.Context, theme *material.Theme, newPasswordView *NewPasswordView, passwordLength string, info Information) {
	elementMargin := layout.Inset{Top: unit.Dp(13), Bottom: unit.Dp(13), Right: unit.Dp(10), Left: unit.Dp(10)}
	btnsMargin := layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20), Right: unit.Dp(10), Left: unit.Dp(10)}
	randomBtnsMargin := layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0), Right: unit.Dp(10), Left: unit.Dp(10)}
	appTextSize := unit.Sp(15)

	layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(60), Right: unit.Dp(60)}.Layout(
		*gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceSides}.Layout(
				gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return elementMargin.Layout(
							gtx,
							func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceSides}.Layout(
									gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										header := material.H3(theme, "New Password")
										header.Font.Typeface = "Verdana, monospace"
										return header.Layout(gtx)
									}),
								)
							},
						)
					},
				),
				horizontalDivider(),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return elementMargin.Layout(
							gtx,
							func(gtx layout.Context) layout.Dimensions {
								label := material.Label(theme, appTextSize, info.text)
								label.Color = info.color
								label.Font.Weight = font.Bold
								return label.Layout(gtx)
							},
						)
					},
				),
				horizontalDivider(),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							return material.H6(theme, "Master Password:").Layout(gtx)
						},
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							inputMasterPassword := material.Editor(theme, newPasswordView.masterPassword, "Enter master Password...")
							inputMasterPassword.TextSize = appTextSize
							inputMasterPassword.SelectionColor = blue
							// border := widget.Border{Color: newPasswordView.borderColor, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}

							// return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(10)).Layout(gtx, inputMasterPassword.Layout)
							// })
						},
					)
				}),
				horizontalDivider(),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							return material.H6(theme, "Service Name:").Layout(gtx)
						},
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							inputMasterPasswordRepeat := material.Editor(theme, newPasswordView.serviceName, "Enter name of service...")
							inputMasterPasswordRepeat.TextSize = appTextSize
							inputMasterPasswordRepeat.SelectionColor = blue

							return layout.UniformInset(unit.Dp(10)).Layout(gtx, inputMasterPasswordRepeat.Layout)
						},
					)
				}),
				horizontalDivider(),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							return material.H6(theme, "Username:").Layout(gtx)
						},
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							inputMasterPasswordRepeat := material.Editor(theme, newPasswordView.username, "Enter username...")
							inputMasterPasswordRepeat.TextSize = appTextSize
							inputMasterPasswordRepeat.SelectionColor = blue

							return layout.UniformInset(unit.Dp(10)).Layout(gtx, inputMasterPasswordRepeat.Layout)
						},
					)
				}),
				horizontalDivider(),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							return material.H6(theme, "Password ["+passwordLength+"]").Layout(gtx)
						},
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							inputMasterPasswordRepeat := material.Editor(theme, newPasswordView.password, "Enter or generate random password...")
							inputMasterPasswordRepeat.TextSize = appTextSize
							inputMasterPasswordRepeat.SelectionColor = blue

							return layout.UniformInset(unit.Dp(10)).Layout(gtx, inputMasterPasswordRepeat.Layout)
						},
					)
				}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return randomBtnsMargin.Layout(
							gtx,
							func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.End, Spacing: layout.SpaceAround}.Layout(
									gtx,
									layout.Flexed(1,
										func(gtx layout.Context) layout.Dimensions {
											return randomBtnsMargin.Layout(
												gtx,
												func(gtx layout.Context) layout.Dimensions {
													confirmBtn := material.Button(theme, newPasswordView.smallRandWidget, "20 - 25")
													confirmBtn.Background = grey_light
													confirmBtn.Color = black
													confirmBtn.TextSize = appTextSize - 5
													confirmBtn.Font.Weight = font.Bold

													return layout.UniformInset(unit.Dp(0)).Layout(gtx, confirmBtn.Layout)
												},
											)
										},
									),
									layout.Flexed(1,
										func(gtx layout.Context) layout.Dimensions {
											return randomBtnsMargin.Layout(
												gtx,
												func(gtx layout.Context) layout.Dimensions {
													confirmBtn := material.Button(theme, newPasswordView.mediumRandWidget, "25 - 40")
													confirmBtn.Background = grey_light
													confirmBtn.Color = black
													confirmBtn.TextSize = appTextSize - 5
													confirmBtn.Font.Weight = font.Bold

													return layout.UniformInset(unit.Dp(0)).Layout(gtx, confirmBtn.Layout)
												},
											)
										},
									),
									layout.Flexed(1,
										func(gtx layout.Context) layout.Dimensions {
											return randomBtnsMargin.Layout(
												gtx,
												func(gtx layout.Context) layout.Dimensions {
													confirmBtn := material.Button(theme, newPasswordView.bigRandWidget, "40-100")
													confirmBtn.Background = grey_light
													confirmBtn.Color = black
													confirmBtn.TextSize = appTextSize - 5
													confirmBtn.Font.Weight = font.Bold

													return layout.UniformInset(unit.Dp(0)).Layout(gtx, confirmBtn.Layout)
												},
											)
										},
									),
									layout.Flexed(1,
										func(gtx layout.Context) layout.Dimensions {
											return randomBtnsMargin.Layout(
												gtx,
												func(gtx layout.Context) layout.Dimensions {
													confirmBtn := material.Button(theme, newPasswordView.specialCharsSwitchWidget, newPasswordView.specialCharsSwitchText)
													confirmBtn.Background = newPasswordView.specialCharsSwitchColor
													confirmBtn.Color = black
													confirmBtn.TextSize = appTextSize - 5
													confirmBtn.Font.Weight = font.Bold

													return layout.UniformInset(unit.Dp(0)).Layout(gtx, confirmBtn.Layout)
												},
											)
										},
									),
								)
							},
						)
					},
				),
				emptyDivider(),
				horizontalDivider(),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return btnsMargin.Layout(
							gtx,
							func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceSides}.Layout(
									gtx,
									layout.Rigid(
										func(gtx layout.Context) layout.Dimensions {
											return elementMargin.Layout(
												gtx,
												func(gtx layout.Context) layout.Dimensions {
													confirmBtn := material.Button(theme, newPasswordView.confirmBtnWidget, "            SAVE            ")
													confirmBtn.Background = purple_light
													confirmBtn.TextSize = appTextSize
													confirmBtn.Font.Weight = font.Normal
													confirmBtn.Color = black
													confirmBtn.Font.Typeface = "Verdana, monospace"

													return layout.UniformInset(unit.Dp(0)).Layout(gtx, confirmBtn.Layout)
												},
											)
										},
									),
									layout.Rigid(
										func(gtx layout.Context) layout.Dimensions {
											return elementMargin.Layout(
												gtx,
												func(gtx layout.Context) layout.Dimensions {
													confirmBtn := material.Button(theme, newPasswordView.showHidWidget, "SHOW/HIDE")
													confirmBtn.Background = grey_light
													confirmBtn.TextSize = appTextSize
													confirmBtn.Font.Weight = font.Normal
													confirmBtn.Color = black
													confirmBtn.Font.Typeface = "Verdana, monospace"

													return confirmBtn.Layout(gtx)
												},
											)
										},
									),
								)
							},
						)
					},
				),
			)
		},
	)
}

func ConfirmPasswordDeletionWidget(gtx *layout.Context, theme *material.Theme, serviceName string, confirm *widget.Clickable, deny *widget.Clickable) {
	var (
		textSize    unit.Sp      = 30
		btnMargin   layout.Inset = layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20), Right: unit.Dp(25), Left: unit.Dp(25)}
		labelMargin layout.Inset = layout.Inset{Top: unit.Dp(25), Bottom: unit.Dp(25), Right: unit.Dp(25), Left: unit.Dp(25)}
	)

	layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides}.Layout(
		*gtx,
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				return labelMargin.Layout(
					gtx,
					func(gtx layout.Context) layout.Dimensions {
						label := material.Label(theme, textSize, "Do you want to delete all information about service "+serviceName+"? This action can't be reversed.")
						label.Color = red
						label.Font.Typeface = "Verdana, monospace"
						return label.Layout(gtx)
					},
				)
			},
		),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(
					gtx,
					layout.Flexed(1,
						func(gtx layout.Context) layout.Dimensions {
							return btnMargin.Layout(
								gtx,
								func(gtx layout.Context) layout.Dimensions {
									confirmBtn := material.Button(theme, confirm, "CONFIRM")
									confirmBtn.Font.Weight = font.Bold
									confirmBtn.Background = red
									confirmBtn.Color = black
									confirmBtn.Font.Typeface = "Verdana, monospace"

									return confirmBtn.Layout(gtx)
								},
							)
						},
					),
					layout.Flexed(1,
						func(gtx layout.Context) layout.Dimensions {
							return btnMargin.Layout(
								gtx,

								func(gtx layout.Context) layout.Dimensions {
									denyBtn := material.Button(theme, deny, "DENY")
									denyBtn.Font.Weight = font.Bold
									denyBtn.Background = grey_light
									denyBtn.Color = black
									denyBtn.Font.Typeface = "Verdana, monospace"

									return denyBtn.Layout(gtx)
								},
							)
						},
					),
				)
			},
		),
	)
}

func ResizeDecryptionWindow(window *app.Window) {
	var (
		maxW unit.Dp = 850
		maxH unit.Dp = 800
	)
	window.Option(app.Decorated(true))
	window.Option(app.MinSize(unit.Dp(maxW), unit.Dp(maxH)))
	window.Option(app.MaxSize(unit.Dp(maxW*2), unit.Dp(maxH)))
	window.Option(app.Size(unit.Dp(maxW), unit.Dp(maxH)))
	window.Option(app.Title(appName))

	return
}

func ManagePasswordDecryptionWidget(gtx *layout.Context, theme *material.Theme, serviceName *string, textCheckMsg *string, authenticate *widget.Clickable, cancel *widget.Clickable, showHideUsername *widget.Clickable, showHidePassword *widget.Clickable, masterPasswordGUI *widget.Editor, usernameGUI *widget.Editor, passwordGUI *widget.Editor, passwordEditorBackgroundColor *color.NRGBA) {
	var (
		appTextSize       unit.Sp      = 20
		btnMargin         layout.Inset = layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20), Right: unit.Dp(25), Left: unit.Dp(25)}
		showHideBtnMargin layout.Inset = layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Right: unit.Dp(20), Left: unit.Dp(0)}
		elementMargin     layout.Inset = layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Right: unit.Dp(20), Left: unit.Dp(20)}
	)

	layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(60), Right: unit.Dp(60)}.Layout(
		*gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides}.Layout(
				gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {

						return layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(20), Left: unit.Dp(0), Right: unit.Dp(0)}.Layout(
							gtx,
							func(gtx layout.Context) layout.Dimensions {

								return elementMargin.Layout(
									gtx,
									func(gtx layout.Context) layout.Dimensions {
										label := material.Label(theme, textSize, "Decrypt information for "+*serviceName)
										label.Color = black
										return label.Layout(gtx)
									},
								)
							},
						)
					},
				),
				horizontalDivider(),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return elementMargin.Layout(
							gtx,
							func(gtx layout.Context) layout.Dimensions {
								masterPasswordLabel := material.H6(theme, "Master Password"+*textCheckMsg)
								if len(*textCheckMsg) == 0 {
									masterPasswordLabel.Color = black
								} else {
									masterPasswordLabel.Color = red
								}

								return masterPasswordLabel.Layout(gtx)
							},
						)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return elementMargin.Layout(
							gtx,
							func(gtx layout.Context) layout.Dimensions {
								inputMasterPasswordRepeat := material.Editor(theme, masterPasswordGUI, "Enter master password and authenticate...")
								inputMasterPasswordRepeat.HintColor = blue
								inputMasterPasswordRepeat.TextSize = appTextSize
								inputMasterPasswordRepeat.SelectionColor = blue

								return layout.UniformInset(unit.Dp(10)).Layout(gtx, inputMasterPasswordRepeat.Layout)
							},
						)
					},
				),
				horizontalDivider(),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							return material.H6(theme, "Username").Layout(gtx)
						},
					)
				}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(
							gtx,
							layout.Flexed(
								1,
								func(gtx layout.Context) layout.Dimensions {
									return elementMargin.Layout(
										gtx,
										func(gtx layout.Context) layout.Dimensions {
											usernameEditor := material.Editor(theme, usernameGUI, "Username will show after authentication...")
											usernameEditor.SelectionColor = blue
											usernameEditor.Font.Typeface = "Verdana, monospace"
											usernameEditor.Font.Weight = font.Medium

											return layout.UniformInset(unit.Dp(10)).Layout(gtx, usernameEditor.Layout)
										},
									)
								},
							),
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {

									return showHideBtnMargin.Layout(
										gtx,
										func(gtx layout.Context) layout.Dimensions {
											showHideBtn := material.Button(theme, showHideUsername, "show")
											showHideBtn.Inset = layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12), Left: unit.Dp(23), Right: unit.Dp(23)}
											showHideBtn.TextSize = appTextSize
											showHideBtn.Background = grey_light
											showHideBtn.Font.Typeface = "Verdana, monospace"
											showHideBtn.Font.Weight = font.Normal
											showHideBtn.Color = black

											return showHideBtn.Layout(gtx)
										},
									)
								},
							),
						)
					},
				),
				horizontalDivider(),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return elementMargin.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							return material.H6(theme, "Password").Layout(gtx)
						},
					)
				}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(
							gtx,
							layout.Flexed(
								1,
								func(gtx layout.Context) layout.Dimensions {

									return elementMargin.Layout(
										gtx,
										func(gtx layout.Context) layout.Dimensions {
											passwordEditor := material.Editor(theme, passwordGUI, "Password will show after authentication...")
											passwordEditor.TextSize = appTextSize
											passwordEditor.SelectionColor = blue
											passwordEditor.Font.Typeface = "Verdana, monospace"
											passwordEditor.Font.Weight = font.Medium

											return layout.UniformInset(unit.Dp(10)).Layout(gtx, passwordEditor.Layout)
										},
									)

								},
							),
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									return showHideBtnMargin.Layout(
										gtx,
										func(gtx layout.Context) layout.Dimensions {
											showHideBtn := material.Button(theme, showHidePassword, "show")
											showHideBtn.Inset = layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12), Left: unit.Dp(23), Right: unit.Dp(23)}
											showHideBtn.TextSize = appTextSize
											showHideBtn.Background = grey_light
											showHideBtn.Font.Typeface = "Verdana, monospace"
											showHideBtn.Font.Weight = font.Medium
											showHideBtn.Color = black

											return layout.UniformInset(unit.Dp(0)).Layout(gtx, showHideBtn.Layout)
										},
									)
								},
							),
						)
					},
				),
				horizontalDivider(),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(
							gtx,
							layout.Flexed(1,
								func(gtx layout.Context) layout.Dimensions {
									return btnMargin.Layout(
										gtx,
										func(gtx layout.Context) layout.Dimensions {
											confirmBtn := material.Button(theme, authenticate, "AUTHENTICATE")
											confirmBtn.TextSize = appTextSize
											confirmBtn.Font.Weight = font.Medium
											confirmBtn.Background = blue
											confirmBtn.Color = black
											confirmBtn.Font.Typeface = "Verdana, monospace"

											return layout.UniformInset(unit.Dp(0)).Layout(gtx, confirmBtn.Layout)
										},
									)
								},
							),
							layout.Flexed(1,
								func(gtx layout.Context) layout.Dimensions {
									return btnMargin.Layout(
										gtx,

										func(gtx layout.Context) layout.Dimensions {
											cancelBtn := material.Button(theme, cancel, "CANCEL")
											cancelBtn.TextSize = appTextSize
											cancelBtn.Font.Weight = font.Medium
											cancelBtn.Background = grey_light
											cancelBtn.Color = black
											cancelBtn.Font.Typeface = "Verdana, monospace"

											return layout.UniformInset(unit.Dp(0)).Layout(gtx, cancelBtn.Layout)
										},
									)
								},
							),
						)
					},
				),
			)
		},
	)
}

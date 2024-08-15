package gui

import (
	"image/color"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	beige    = color.NRGBA{R: 247, G: 239, B: 229, A: 255}
	black    = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	blue     = color.NRGBA{R: 150, G: 201, B: 244, A: 255}
	green    = color.NRGBA{R: 100, G: 196, B: 166, A: 255}
	grey     = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	orange   = color.NRGBA{R: 235, G: 178, B: 10, A: 100}
	purple   = color.NRGBA{R: 161, G: 100, B: 196, A: 255}
	red      = color.NRGBA{R: 238, G: 78, B: 78, A: 255}
	textSize = unit.Sp(30)
)

const alphabet = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789" + "!#$%&'()*+,-./:;<=>?@"

func ResizeWindowInfo(window *app.Window) {
	window.Option(app.Decorated(true))
	window.Option(app.MinSize(unit.Dp(300), unit.Dp(300)))
	window.Option(app.MaxSize(unit.Dp(450), unit.Dp(450)))
	window.Option(app.Size(unit.Dp(450), unit.Dp(450)))
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

				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return elementMargin.Layout(
							gtx,
							func(gtx layout.Context) layout.Dimensions {
								label := material.Label(theme, unit.Sp(20), "Please enter your master password. The longer your password, the better it will protect your sensitive information. Remember, this master password will encrypt and secure all your other passwords. If you forget it, there is no way to recover it—your access will be permanently lost.")
								label.Color = red
								label.Font.Weight = font.Bold
								return label.Layout(gtx)
							},
						)
					},
				),
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
							inputMasterPassword := material.Editor(theme, initialSetup.passwordInput, "Password")
							inputMasterPassword.TextSize = unit.Sp(20)
							inputMasterPassword.SelectionColor = blue
							border := widget.Border{Color: initialSetup.borderColor, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}

							return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.UniformInset(unit.Dp(10)).Layout(gtx, inputMasterPassword.Layout)
							})
						},
					)
				}),
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
							inputMasterPasswordRepeat := material.Editor(theme, initialSetup.passwordInputRepeat, "Password")
							inputMasterPasswordRepeat.TextSize = unit.Sp(20)
							inputMasterPasswordRepeat.SelectionColor = blue
							border := widget.Border{Color: initialSetup.borderColor, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}

							return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.UniformInset(unit.Dp(10)).Layout(gtx, inputMasterPasswordRepeat.Layout)
							})
						},
					)
				}),
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
													confirmBtn := material.Button(theme, initialSetup.confirmBtnWidget, "CONFIRM")
													confirmBtn.Background = green
													confirmBtn.TextSize = unit.Sp(25)
													confirmBtn.Font.Weight = font.Bold
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
													confirmBtn := material.Button(theme, initialSetup.showHidWidget, "SHOW/HIDE")
													confirmBtn.Background = color.NRGBA{R: 30, G: 30, B: 30, A: 255}
													confirmBtn.TextSize = unit.Sp(25)
													confirmBtn.Font.Weight = font.Bold
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
					loader.Color = green
					return loader.Layout(gtx)
				}),
			)
		},
	)
}

func ResizeWindowPasswordEntriesList(window *app.Window) {
	window.Option(app.Decorated(true))
	window.Option(app.MinSize(unit.Dp(300), unit.Dp(300)))
	window.Option(app.MaxSize(unit.Dp(2000), unit.Dp(2000)))
	window.Option(app.Size(unit.Dp(450), unit.Dp(800)))
}

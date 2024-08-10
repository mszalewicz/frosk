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
	black    = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	green    = color.NRGBA{R: 100, G: 196, B: 166, A: 255}
	orange   = color.NRGBA{R: 235, G: 178, B: 10, A: 100}
	red      = color.NRGBA{R: 238, G: 78, B: 78, A: 255}
	purple   = color.NRGBA{R: 161, G: 100, B: 196, A: 255}
	beige    = color.NRGBA{R: 247, G: 239, B: 229, A: 255}
	grey     = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	textSize = unit.Sp(30)
)

func ResizeWindowInfo(window *app.Window) {
	window.Option(app.Decorated(true))
	window.Option(app.MinSize(unit.Dp(300), unit.Dp(300)))
	window.Option(app.MaxSize(unit.Dp(450), unit.Dp(450)))
	window.Option(app.Size(unit.Dp(450), unit.Dp(450)))
}

func ResizeWindowInitialSetup(window *app.Window) {
	window.Option(app.Decorated(true))
	window.Option(app.MinSize(unit.Dp(300), unit.Dp(300)))
	window.Option(app.MaxSize(unit.Dp(2000), unit.Dp(2000)))
	window.Option(app.Size(unit.Dp(1_030), unit.Dp(1_000)))
}

func InfoWindowWidget(gtx *layout.Context, theme *material.Theme, returnBtnWidget *widget.Clickable, list *widget.List, text *string) {

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
								return material.Label(theme, unit.Sp(18), *text).Layout(gtx)
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

func InitialSetupWidget(gtx *layout.Context, theme *material.Theme) {
	elementMargin := layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(30)}

	layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20), Left: unit.Dp(60), Right: unit.Dp(60)}.Layout(
		*gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(
				gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return elementMargin.Layout(
							gtx,
							func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceSides}.Layout(
									gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions { return material.H1(theme, "Initial Setup").Layout(gtx) }),
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
				// layout.Flexed(
				// 	1,
				// 	func(gtx layout.Context) layout.Dimensions {
				// 		material.edi
				// 	},
				// ),
			)
		},
	)
}

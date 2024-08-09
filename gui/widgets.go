package gui

import (
	"image/color"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/system"
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
	textSize = unit.Sp(30)
)

func ResizeWindowInfo(window *app.Window) {
	window.Option(app.Decorated(true))
	window.Option(app.MinSize(unit.Dp(300), unit.Dp(300)))
	window.Option(app.MaxSize(unit.Dp(450), unit.Dp(450)))
	window.Option(app.Size(unit.Dp(450), unit.Dp(450)))
	window.Perform(system.ActionCenter)
}

func TransformIntoInfoWindow(gtx *layout.Context, theme *material.Theme, returnBtnWidget *widget.Clickable, list *widget.List, text *string) {

	layout.Inset{Top: 20, Bottom: 20, Left: 20, Right: 20}.Layout(
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
								return material.Label(theme, unit.Sp(30), *text).Layout(gtx)
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

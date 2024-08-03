package gui

import (
	"image/color"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func HandleMainWindow(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops

	var mainButtonWidget widget.Clickable
	var margins = layout.Inset{Top: unit.Dp(25), Bottom: unit.Dp(25), Left: unit.Dp(35), Right: unit.Dp(35)}

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(
				gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								mainBtn := material.Button(theme, &mainButtonWidget, "test")
								mainBtn.Background = color.NRGBA{R: 70, G: 72, B: 72, A: 255}
								mainBtn.TextSize = unit.Sp(42)
								return mainBtn.Layout(gtx)
							},
						)
					},
				),
				layout.Rigid(layout.Spacer{Height: unit.Dp(25)}.Layout),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						title := material.H1(theme, "Test")
						title.Color = color.NRGBA{R: 127, G: 20, B: 120, A: 255}
						title.Alignment = text.Middle
						return title.Layout(gtx)
					},
				),
				layout.Rigid(layout.Spacer{Height: unit.Dp(25)}.Layout),
			)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}

// func HandleErrorWindow(window *app.Window) error {
// 	theme := material.NewTheme()
// 	var ops op.Ops
// 	for {
// 		switch e := window.Event().(type) {
// 		case app.DestroyEvent:
// 			return e.Err
// 		case app.FrameEvent:
// 			// This graphics context is used for managing the rendering state.
// 			gtx := app.NewContext(&ops, e)

// 			// Define an large label with an appropriate text:
// 			title := material.H1(theme, "Test")

// 			// Change the color of the label.
// 			maroon := color.NRGBA{R: 127, G: 20, B: 120, A: 255}
// 			title.Color = maroon

// 			// Change the position of the label.
// 			title.Alignment = text.Middle

// 			// Draw the label to the graphics context.
// 			title.Layout(gtx)

// 			// Pass the drawing operations to the GPU.
// 			e.Frame(gtx.Ops)
// 		}
// 	}
// }

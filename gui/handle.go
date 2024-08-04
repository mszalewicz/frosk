package gui

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func createPasswordEntryLineComponents(serviceName string, theme *material.Theme) []layout.FlexChild {
	const buttonSize = 12

	var openBtnWidget widget.Clickable
	openBtn := material.Button(theme, &openBtnWidget, "OPEN")
	openBtn.Background = color.NRGBA{R: 67, G: 168, B: 84, A: 255}
	openBtn.TextSize = unit.Sp(buttonSize)

	var deleteBtnWidget widget.Clickable
	deleteBtn := material.Button(theme, &deleteBtnWidget, "DELETE")
	deleteBtn.Background = color.NRGBA{R: 235, G: 64, B: 52, A: 255}
	deleteBtn.TextSize = unit.Sp(buttonSize)

	var btnMargin = layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(0)}
	var labelMargin = layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(0)}

	serviceFlexChild := layout.Flexed(
		1,
		func(gtx layout.Context) layout.Dimensions {
			return labelMargin.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {
					return material.Label(theme, unit.Sp(25), serviceName).Layout(gtx)
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

func constructPasswordEntriesList(passwordEntries [][]layout.FlexChild, passwordEntriesList layout.List, margin layout.Inset) layout.FlexChild {
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
								horizontalSpacer(),
							)
						},
					)
				},
			)
		},
	)
}

func horizontalSpacer() layout.FlexChild {
	return layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			height := unit.Dp(1)
			line := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Dp(height))
			paint.FillShape(gtx.Ops, color.NRGBA{A: 40}, clip.Rect(line).Op())
			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Dp(height))}
		},
	)
}

func HandleMainWindow(window *app.Window) error {
	theme := material.NewTheme()
	initialRender := true
	var ops op.Ops
	var newPasswordEntryWidget widget.Clickable
	var margin = layout.Inset{Top: unit.Dp(15), Bottom: unit.Dp(15), Left: unit.Dp(15), Right: unit.Dp(15)}

	testServices := []string{"google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank", "google", "email", "facebook", "twitter", "bank"}
	fmt.Println(len(testServices))

	passwordEntriesList := layout.List{Axis: layout.Vertical}
	passwordEntries := [][]layout.FlexChild{}

	for _, serviceName := range testServices {
		passwordEntries = append(passwordEntries, createPasswordEntryLineComponents(serviceName, theme))
	}

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
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

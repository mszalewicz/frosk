package gui

import (
	"fmt"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// In some cases Gio UI do not calculate sizes of elements correctly.
// For example when centering windows via build-in functionality.
// This helper function prompts window to repaint.
// Signal has to come from outside of the window goroutine.
func refresh_window(window *app.Window) {
	go func() {
		for range 3 {
			time.Sleep(time.Second / 20)
			window.Invalidate()
		}
		return
	}()
}

// Pop-out information window.
func Render_Information_Window(
	gtx *layout.Context,
	theme *material.Theme,
	returnBtnWidget *widget.Clickable,
	list *widget.List,
	text string,
) {

	returnButtonElement, returnButtonWidget := CreateDefaultButton(
		theme,
		"Return",
		32,
		DEFAULT_PADDING,
		DEFAULT_PADDING,
		DEFAULT_PADDING,
		DEFAULT_PADDING,
	)

	returnButtonElement.Button = returnBtnWidget

	renderList := func(gtx layout.Context) layout.Dimensions {
		return list.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
			return material.Label(theme, unit.Sp(DEFAULT_TEXT_SIZE*2), text).Layout(gtx)
		})
	}

	// Define the uniform spacing
	margins := layout.Inset{
		Top:    unit.Dp(BIG_PADDING),
		Bottom: unit.Dp(BIG_PADDING),
		Left:   unit.Dp(BIG_PADDING),
		Right:  unit.Dp(BIG_PADDING),
	}

	// Execute the layout
	margins.Layout(*gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.End,
		}.Layout(gtx,
			layout.Flexed(1, renderList),     // Content takes up remaining space
			layout.Rigid(returnButtonWidget), // Button stays at the bottom
		)
	})
}

func Error_Window(ops *op.Ops, window *app.Window, theme *material.Theme, errorMsg string) error {
	window.Option(app.Decorated(false))
	window.Option(app.MinSize(unit.Dp(300), unit.Dp(300)))
	window.Option(app.MaxSize(unit.Dp(450), unit.Dp(450)))
	window.Option(app.Size(unit.Dp(450), unit.Dp(450)))
	window.Option(app.Title(APP_NAME))

	// refresh_window(window)

	errConfirmWidget := new(widget.Clickable)
	errListContainer := &widget.List{
		List: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Start,
		},
	}

	// Render Loop
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(1)
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(ops, e)

			paint.Fill(ops, GREY_LIGHT)

			if errConfirmWidget.Clicked(gtx) {
				os.Exit(1)
				fmt.Println("exit")
			}

			Render_Information_Window(&gtx, theme, errConfirmWidget, errListContainer, errorMsg)

			e.Frame(gtx.Ops)
		}
	}
}

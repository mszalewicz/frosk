package gui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// SplitView is a container that show two widgets side by side.
// Ratio is a measure of how much space percentae wise left widget will take.
// Right side widget will automatically fill rest of avaiable space.
type SplitView struct{}

func (s SplitView) Layout(
	gtx layout.Context,
	left layout.Widget,
	right layout.Widget,
	ratio float32,
) layout.Dimensions {

	if ratio < 0.0 || ratio > 1.0 {
		ratio = 0.5
	}

	leftsize := int(float32(gtx.Constraints.Min.X) * ratio)
	rightsize := gtx.Constraints.Min.X - leftsize

	{
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(leftsize, gtx.Constraints.Max.Y))
		left(gtx)

	}

	{
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(rightsize, gtx.Constraints.Max.Y))
		transpose := op.Offset(image.Pt(leftsize, 0)).Push(gtx.Ops)
		right(gtx)
		transpose.Pop()
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// Create buttom with custom theming and padding - used for special case buttons, e.x. "DELETE ALL" with red background
func CreateUniqueButton(theme *material.Theme, label string, background color.NRGBA, text_color color.NRGBA, font_size int, top_inset int, bottom_inset int, left_inset int, right_inset int) (*material.ButtonStyle, layout.Widget) {

	if font_size == 0 {
		font_size = int(theme.TextSize)
	}

	// 1. Create associated clickable widget
	var buttonWidget widget.Clickable

	// 2. Initialize the material button
	btn := material.Button(theme, &buttonWidget, label)
	btn.Background = background
	btn.Color = text_color
	btn.TextSize = unit.Sp(font_size)

	// 4. Return btn and associated padded layout
	return &btn, func(gtx layout.Context) layout.Dimensions {

		// 2. Define the Inset (the "outer" padding)
		inset := layout.Inset{
			Top:    unit.Dp(top_inset),
			Bottom: unit.Dp(bottom_inset),
			Left:   unit.Dp(left_inset),
			Right:  unit.Dp(right_inset),
		}

		// 3. Return the layout wrapped in the inset
		return inset.Layout(gtx, btn.Layout)
	}
}

// Creates button with default theme and custom padding
func CreateDefaultButton(theme *material.Theme, label string, font_size int, top_inset int, bottom_inset int, left_inset int, right_inset int) (*material.ButtonStyle, layout.Widget) {

	if font_size == 0 {
		font_size = int(theme.TextSize)
	}

	// 1. Create associated clickable widget
	var buttonWidget widget.Clickable

	// 2. Initialize the material button
	btn := material.Button(theme, &buttonWidget, label)
	btn.Background = theme.Bg
	btn.Color = theme.Fg
	btn.TextSize = unit.Sp(font_size)

	// 5. Return button and associated padded layout
	return &btn, func(gtx layout.Context) layout.Dimensions {

		// 3. Define the Inset (the "outer" padding)
		inset := layout.Inset{
			Top:    unit.Dp(top_inset),
			Bottom: unit.Dp(bottom_inset),
			Left:   unit.Dp(left_inset),
			Right:  unit.Dp(right_inset),
		}

		// 4. Return the layout wrapped in the inset
		return inset.Layout(gtx, btn.Layout)
	}
}

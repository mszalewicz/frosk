package gui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

// SplitView is a container that show two widgets side by side.
// Ratio is a measure of how much space percentage wise left widget will take.
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

	leftsize := int( float32(gtx.Constraints.Min.X) * ratio )
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

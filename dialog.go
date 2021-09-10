package razgio

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Dialog struct {
	Message string
	ButtonBar
	Width  unit.Value
	Height unit.Value
	Inset  unit.Value
}

func NewDialog(message string, buttons ...string) Dialog {
	return Dialog{
		Message:   message,
		ButtonBar: NewButtonBar(buttons...),
		Width:     unit.Dp(200),
		Inset:     unit.Dp(16),
	}
}

func (d *Dialog) Layout(gtx C, th *material.Theme) D {
	return layout.Center.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Px(d.Width)
		gtx.Constraints.Min.Y = gtx.Px(d.Height)

		macro := op.Record(gtx.Ops)
		dims := layout.Center.Layout(gtx, func(gtx C) D {
			return layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					in := layout.UniformInset(d.Inset)
					return in.Layout(gtx, material.Body1(th, d.Message).Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return d.ButtonBar.Layout(gtx, th)
				}),
			)
		})
		call := macro.Stop()

		pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
		pointer.InputOp{Tag: d, Types: pointer.Press}.Add(gtx.Ops)

		rr := float32(gtx.Px(unit.Dp(8)))
		clip.UniformRRect(f32.Rectangle{Max: f32.Point{
			X: float32(dims.Size.X),
			Y: float32(dims.Size.Y),
		}}, rr).Add(gtx.Ops)
		paint.Fill(gtx.Ops, th.Bg)

		call.Add(gtx.Ops)

		return dims
	})
}

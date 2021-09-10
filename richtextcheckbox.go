package razgio

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/richtext"
)

// NOTE: most of the code is copied from gio

type RichTextCheckBoxStyle struct {
	richtextCheckable
	CheckBox *widget.Bool
}

func RichTextCheckBox(th *material.Theme, checkBox *widget.Bool, label ...richtext.SpanStyle) RichTextCheckBoxStyle {
	return RichTextCheckBoxStyle{
		CheckBox: checkBox,
		richtextCheckable: richtextCheckable{
			Label:              label,
			IconColor:          th.Palette.ContrastBg,
			Size:               unit.Dp(26),
			shaper:             th.Shaper,
			checkedStateIcon:   th.Icon.CheckBoxChecked,
			uncheckedStateIcon: th.Icon.CheckBoxUnchecked,
		},
	}
}

func (c RichTextCheckBoxStyle) Layout(gtx layout.Context) layout.Dimensions {
	dims := c.layout(gtx, c.CheckBox.Value, c.CheckBox.Hovered())
	gtx.Constraints.Min = dims.Size
	c.CheckBox.Layout(gtx)
	return dims
}

type richtextCheckable struct {
	Label              []richtext.SpanStyle
	IconColor          color.NRGBA
	Size               unit.Value
	shaper             text.Shaper
	checkedStateIcon   *widget.Icon
	uncheckedStateIcon *widget.Icon
}

func (c *richtextCheckable) layout(gtx layout.Context, checked, hovered bool) layout.Dimensions {
	var icon *widget.Icon
	if checked {
		icon = c.checkedStateIcon
	} else {
		icon = c.uncheckedStateIcon
	}

	dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Px(c.Size) * 4 / 3
					dims := layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
					if !hovered {
						return dims
					}

					background := mulAlpha(c.IconColor, 70)

					radius := float32(size) / 2
					paint.FillShape(gtx.Ops, background,
						clip.Circle{
							Center: f32.Point{X: radius, Y: radius},
							Radius: radius,
						}.Op(gtx.Ops))

					return dims
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						size := gtx.Px(c.Size)
						col := c.IconColor
						if gtx.Queue == nil {
							col = disabled(col)
						}
						gtx.Constraints.Min = image.Point{X: size}
						icon.Layout(gtx, col)
						return layout.Dimensions{
							Size: image.Point{X: size, Y: size},
						}
					})
				}),
			)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				var state richtext.InteractiveText
				return richtext.Text(&state, c.shaper, c.Label...).Layout(gtx)
			})
		}),
	)
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	return dims
}

func mulAlpha(c color.NRGBA, alpha uint8) color.NRGBA {
	c.A = uint8(uint32(c.A) * uint32(alpha) / 0xFF)
	return c
}

func disabled(c color.NRGBA) (d color.NRGBA) {
	const r = 80 // blend ratio
	lum := approxLuminance(c)
	return color.NRGBA{
		R: byte((int(c.R)*r + int(lum)*(256-r)) / 256),
		G: byte((int(c.G)*r + int(lum)*(256-r)) / 256),
		B: byte((int(c.B)*r + int(lum)*(256-r)) / 256),
		A: byte(int(c.A) * (128 + 32) / 256),
	}
}

func approxLuminance(c color.NRGBA) byte {
	const (
		r = 13933 // 0.2126 * 256 * 256
		g = 46871 // 0.7152 * 256 * 256
		b = 4732  // 0.0722 * 256 * 256
		t = r + g + b
	)
	return byte((r*int(c.R) + g*int(c.G) + b*int(c.B)) / t)
}

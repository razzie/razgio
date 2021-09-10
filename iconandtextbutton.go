package razgio

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type IconAndTextButtonStyle struct {
	material.ButtonStyle
	Icon   *widget.Icon
	shaper text.Shaper
}

func IconAndTextButton(th *material.Theme, btn *widget.Clickable, icon *widget.Icon, text string) IconAndTextButtonStyle {
	return IconAndTextButtonStyle{
		ButtonStyle: material.Button(th, btn, text),
		Icon:        icon,
		shaper:      th.Shaper,
	}
}

func (b IconAndTextButtonStyle) Layout(gtx C) D {
	return material.ButtonLayoutStyle{
		Background:   b.Background,
		CornerRadius: b.CornerRadius,
		Button:       b.Button,
	}.Layout(gtx, func(gtx C) D {
		return b.ButtonStyle.Inset.Layout(gtx, func(gtx C) D {
			if len(b.Text) == 0 {
				size := gtx.Px(b.TextSize)
				gtx.Constraints = layout.Exact(image.Pt(size, size))
				return b.Icon.Layout(gtx, b.Color)
			}

			iconAndLabel := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
			textIconSpacer := unit.Dp(5)

			layIcon := layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: textIconSpacer}.Layout(gtx, func(gtx C) D {
					size := gtx.Px(b.TextSize)
					gtx.Constraints = layout.Exact(image.Pt(size, size))
					return b.Icon.Layout(gtx, b.Color)
				})
			})

			layLabel := layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx C) D {
					paint.ColorOp{Color: b.Color}.Add(gtx.Ops)
					return widget.Label{Alignment: text.Middle}.Layout(gtx, b.shaper, b.Font, b.TextSize, b.Text)
				})
			})

			return iconAndLabel.Layout(gtx, layIcon, layLabel)
		})
	})
}

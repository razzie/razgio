package razgio

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ButtonBar struct {
	list layout.List
	in   layout.Inset
	btns []buttonBarItem
}

type buttonBarItem struct {
	widget.Clickable
	icon *widget.Icon
	text string
}

func NewButtonBar(btns ...string) ButtonBar {
	bb := ButtonBar{
		list: layout.List{Axis: layout.Horizontal},
		in:   layout.UniformInset(unit.Dp(5)),
		btns: make([]buttonBarItem, len(btns)),
	}
	for i, btn := range btns {
		bb.btns[i].text = btn
	}
	return bb
}

func (bb *ButtonBar) SetButtonIcon(idx int, icon *widget.Icon) {
	bb.btns[idx].icon = icon
}

func (bb *ButtonBar) Clicked(idx int) bool {
	return bb.btns[idx].Clicked()
}

func (bb *ButtonBar) Layout(gtx C, th *material.Theme) D {
	return bb.list.Layout(gtx, len(bb.btns), func(gtx C, idx int) D {
		return bb.in.Layout(gtx, func(gtx C) D {
			return bb.btns[idx].Layout(gtx, th)
		})
	})
}

func (bbi *buttonBarItem) Layout(gtx C, th *material.Theme) D {
	if bbi.icon == nil {
		return material.Button(th, &bbi.Clickable, bbi.text).Layout(gtx)
	} else {
		return IconAndTextButton(th, &bbi.Clickable, bbi.icon, bbi.text).Layout(gtx)
	}
}

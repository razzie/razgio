package razgio

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Tabs struct {
	theme    *material.Theme
	list     layout.List
	tabs     []tab
	selected int
	slider   Slider
	onSelect func(int)
}

type tab struct {
	btn     widget.Clickable
	title   string
	content layout.Widget
}

func NewTabs(th *material.Theme) *Tabs {
	return &Tabs{
		theme: th,
	}
}

func (tabs *Tabs) SetSelectFunc(onSelect func(int)) {
	tabs.onSelect = onSelect
}

func (tabs *Tabs) AddTab(title string, content layout.Widget) int {
	tabs.tabs = append(tabs.tabs, tab{title: title, content: content})
	return len(tabs.tabs) - 1
}

func (tabs *Tabs) SwitchToTab(tabIdx int) {
	if tabs.selected != tabIdx {
		if tabs.selected < tabIdx {
			tabs.slider.PushLeft()
		} else if tabs.selected > tabIdx {
			tabs.slider.PushRight()
		}
		tabs.selected = tabIdx
		if tabs.onSelect != nil {
			tabs.onSelect(tabIdx)
		}
	}
}

func (tabs *Tabs) Layout(gtx C) D {
	gtx.Constraints.Min = gtx.Constraints.Max
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return tabs.list.Layout(gtx, len(tabs.tabs), func(gtx C, tabIdx int) D {
				if len(tabs.tabs[tabIdx].title) == 0 {
					return D{}
				}
				t := &tabs.tabs[tabIdx]
				if t.btn.Clicked() {
					tabs.SwitchToTab(tabIdx)
				}
				var tabWidth int
				return layout.Stack{Alignment: layout.S}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						dims := material.Clickable(gtx, &t.btn, func(gtx C) D {
							return layout.UniformInset(unit.Sp(12)).Layout(gtx,
								material.H6(tabs.theme, t.title).Layout,
							)
						})
						tabWidth = dims.Size.X
						return dims
					}),
					layout.Stacked(func(gtx C) D {
						if tabs.selected != tabIdx {
							return layout.Dimensions{}
						}
						tabHeight := gtx.Px(unit.Dp(4))
						tabRect := f32.Rect(0, 0, float32(tabWidth), float32(tabHeight))
						clip.UniformRRect(tabRect, float32(tabHeight)/2).Add(gtx.Ops)
						paint.Fill(gtx.Ops, tabs.theme.ContrastBg)
						//paint.FillShape(gtx.Ops, tabs.theme.Palette.ContrastBg, clip.Rect(tabRect).Op())
						return layout.Dimensions{
							Size: image.Point{X: tabWidth, Y: tabHeight},
						}
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return tabs.slider.Layout(gtx, func(gtx C) D {
				w := tabs.tabs[tabs.selected].content
				if w != nil {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, w)
				}
				fill(gtx, dynamicColor(tabs.selected), dynamicColor(tabs.selected+1))
				return layout.Center.Layout(gtx,
					material.H1(tabs.theme, fmt.Sprintf("Tab content #%d", tabs.selected+1)).Layout,
				)
			})
		}),
	)
}

func fill(gtx layout.Context, col1, col2 color.NRGBA) {
	dr := image.Rectangle{Max: gtx.Constraints.Min}
	col2.R = byte(float32(col2.R))
	col2.G = byte(float32(col2.G))
	col2.B = byte(float32(col2.B))
	paint.LinearGradientOp{
		Stop1:  f32.Pt(float32(dr.Min.X), 0),
		Stop2:  f32.Pt(float32(dr.Max.X), 0),
		Color1: col1,
		Color2: col2,
	}.Add(gtx.Ops)
	defer op.Save(gtx.Ops).Load()
	clip.Rect(dr).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func dynamicColor(i int) color.NRGBA {
	sn, cs := math.Sincos(float64(i) * math.Phi)
	return color.NRGBA{
		R: 0xA0 + byte(0x30*sn),
		G: 0xA0 + byte(0x30*cs),
		B: 0xD0,
		A: 0xFF,
	}
}

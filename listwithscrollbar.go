package razgio

import (
	"image"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type ListWithScrollbar struct {
	layout.List
	component.VisibilityAnimation
	FinalAlpha         int
	VisibilityDuration time.Duration

	activated  time.Time
	fakeOps    op.Ops
	wSizes     []int
	max        int
	contentLen int
	offset     int
}

func NewListWithScrollbar() ListWithScrollbar {
	return ListWithScrollbar{
		List: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Start,
		},
		VisibilityAnimation: component.VisibilityAnimation{
			State:    component.Invisible,
			Duration: time.Second,
		},
		FinalAlpha:         128,
		VisibilityDuration: time.Second,
	}
}

func (l *ListWithScrollbar) FitsScreen() bool {
	return l.contentLen <= l.max
}

func (l *ListWithScrollbar) update(gtx C, len int, w layout.ListElement) {
	fakeGtx := gtx
	fakeGtx.Ops = &l.fakeOps
	fakeGtx.Reset()
	fakeGtx.Constraints.Min = image.Pt(0, 0)

	l.wSizes = l.wSizes[:0]
	contentLen := 0
	for i := 0; i < len; i++ {
		size := l.Axis.Convert(w(fakeGtx, i).Size).X
		l.wSizes = append(l.wSizes, size)
		contentLen += size
	}
	offset := l.Position.Offset
	if l.Position.First >= len {
		l.Position.First = len - 1
	}
	for i := 0; i < l.Position.First; i++ {
		offset += l.wSizes[i]
	}

	max := l.Axis.Convert(gtx.Constraints.Max).X
	if contentLen < max {
		l.Disappear(gtx.Now)
	} else {
		if l.max != max || l.offset != offset || l.contentLen != contentLen {
			l.Appear(gtx.Now)
			l.activated = gtx.Now
		} else if l.activated.Add(l.VisibilityDuration).Before(gtx.Now) {
			l.Disappear(gtx.Now)
		}
	}
	if l.Visible() {
		op.InvalidateOp{At: gtx.Now.Add(l.VisibilityDuration)}.Add(gtx.Ops)
	}

	l.max = max
	l.contentLen = contentLen
	l.offset = offset
}

func (l *ListWithScrollbar) Layout(gtx C, th *material.Theme, len int, w layout.ListElement) D {
	l.update(gtx, len, w)

	var stack layout.Stack
	if l.Axis == layout.Vertical {
		stack.Alignment = layout.NE
	} else {
		stack.Alignment = layout.SW
	}
	return stack.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return l.List.Layout(gtx, len, w)
		}),
		layout.Stacked(func(gtx C) D {
			scale := float32(l.max) / float32(l.contentLen)
			if scale > 1 {
				scale = 1
			}
			scrollbarThickness := float32(gtx.Px(unit.Dp(8)))
			scrollbarStart := float32(l.offset) * scale
			scrollbarLen := float32(l.max) * scale
			var scrollbar f32.Rectangle
			if l.Axis == layout.Vertical {
				scrollbar = f32.Rect(
					0,
					scrollbarStart,
					scrollbarThickness,
					scrollbarStart+scrollbarLen,
				)
			} else {
				scrollbar = f32.Rect(
					scrollbarStart,
					0,
					scrollbarStart+scrollbarLen,
					scrollbarThickness,
				)
			}
			rr := scrollbarThickness / 2
			clip.UniformRRect(scrollbar, rr).Add(gtx.Ops)
			alpha := uint8(float32(l.FinalAlpha) * l.Revealed(gtx))
			paint.Fill(gtx.Ops, component.WithAlpha(th.ContrastBg, alpha))
			return D{}
		}),
	)
}

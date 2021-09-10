package razgio

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/richtext"
)

type Tree struct {
	widget.Bool
	Label    []TreeLabel
	Children []*Tree
}

type TreeLabel struct {
	Text      string
	Highlight bool
}

func NewTree(label ...TreeLabel) Tree {
	return Tree{Label: label}
}

func (tree *Tree) AddChild(label ...TreeLabel) *Tree {
	child := &Tree{Label: label}
	tree.Children = append(tree.Children, child)
	return child
}

func (tree *Tree) ClearChildren() {
	tree.Children = nil
}

func (tree *Tree) Layout(gtx C, th *material.Theme) D {
	return tree.layout(gtx, convertToTreeTheme(th))
}

func (tree *Tree) layout(gtx C, th *material.Theme) D {
	if len(tree.Children) == 0 {
		gtx.Queue = nil
	}
	spans := make([]richtext.SpanStyle, len(tree.Label))
	for i, span := range tree.Label {
		color := th.Fg
		if span.Highlight {
			color = th.ContrastBg
		}
		spans[i] = richtext.SpanStyle{
			Size:    th.TextSize.Scale(14.0 / 16.0),
			Color:   color,
			Content: span.Text,
		}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(RichTextCheckBox(th, &tree.Bool, spans...).Layout),
		layout.Rigid(func(gtx C) D {
			if !tree.Value {
				return D{}
			}
			children := make([]layout.FlexChild, len(tree.Children))
			for i, child := range tree.Children {
				childCopy := child
				children[i] = layout.Rigid(func(gtx C) D {
					return childCopy.layout(gtx, th)
				})
			}
			return layout.Inset{Left: th.TextSize.Scale(2)}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
			})
		}),
	)
}

func convertToTreeTheme(th *material.Theme) *material.Theme {
	clone := *th
	clone.Icon.CheckBoxChecked = GetIcons().ArrowDown
	clone.Icon.CheckBoxUnchecked = GetIcons().ArrowRight
	return &clone
}

package razgio

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type MessageBox struct {
	Dialog
	theme  *material.Theme
	okFunc func()
}

func NewMessageBox(th *material.Theme, message string, okFunc func()) layout.Widget {
	mbox := &MessageBox{
		Dialog: NewDialog(message, "OK"),
		theme:  th,
		okFunc: okFunc,
	}
	return mbox.Layout
}

func (mbox *MessageBox) Layout(gtx C) D {
	if mbox.Clicked(0) && mbox.okFunc != nil {
		mbox.okFunc()
	}
	return mbox.Dialog.Layout(gtx, mbox.theme)
}

type YesNoMessageBox struct {
	Dialog
	theme     *material.Theme
	yesNoFunc func(bool)
}

func NewYesNoMessageBox(th *material.Theme, message string, yesNoFunc func(bool)) layout.Widget {
	mbox := &YesNoMessageBox{
		Dialog:    NewDialog(message, "Yes", "No"),
		theme:     th,
		yesNoFunc: yesNoFunc,
	}
	return mbox.Layout
}

func (mbox *YesNoMessageBox) Layout(gtx C) D {
	if mbox.yesNoFunc != nil {
		switch {
		case mbox.Clicked(0):
			mbox.yesNoFunc(true)
		case mbox.Clicked(1):
			mbox.yesNoFunc(false)
		}
	}
	return mbox.Dialog.Layout(gtx, mbox.theme)
}

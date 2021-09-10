package razgio

import (
	"time"

	"gioui.org/layout"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type PageHandler struct {
	theme *material.Theme
	tabs  *Tabs
	pages []Page
	modal component.ModalState
}

func NewPageHandler(th *material.Theme) *PageHandler {
	ph := &PageHandler{
		theme: th,
		tabs:  NewTabs(th),
	}
	ph.tabs.SetSelectFunc(ph.selectPage)
	ph.modal.VisibilityAnimation.Duration = time.Millisecond * 250
	ph.modal.VisibilityAnimation.State = component.Invisible
	return ph
}

func (ph *PageHandler) AddPage(page Page) int {
	ph.pages = append(ph.pages, page)
	return ph.tabs.AddTab(page.GetName(), page.Layout)
}

func (ph *PageHandler) SwitchToPage(idx int) {
	ph.tabs.SwitchToTab(idx)
}

func (ph *PageHandler) ModalMessage(msg string) {
	mbox := NewMessageBox(ph.theme, msg, func() {
		ph.modal.Disappear(time.Now())
	})
	ph.modal.Show(time.Now(), mbox)
}

func (ph *PageHandler) ModalYesNo(msg string, yesFunc func()) {
	mbox := NewYesNoMessageBox(ph.theme, msg, func(yes bool) {
		if yes {
			yesFunc()
		}
		ph.modal.Disappear(time.Now())
	})
	ph.modal.Show(time.Now(), mbox)
}

func (ph *PageHandler) Layout(gtx C) D {
	gtx.Constraints.Min = gtx.Constraints.Max
	return layout.Stack{Alignment: layout.NW}.Layout(gtx,
		layout.Stacked(ph.tabs.Layout),
		layout.Stacked(component.Modal(ph.theme, &ph.modal).Layout),
	)
}

func (ph *PageHandler) selectPage(idx int) {
	ph.pages[idx].Select()
}

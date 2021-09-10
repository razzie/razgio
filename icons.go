package razgio

import (
	"sync"

	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	iconsOnce sync.Once
	icns      *Icons
)

type Icons struct {
	ContentAdd            *widget.Icon
	ActionDelete          *widget.Icon
	CheckBoxBlank         *widget.Icon
	CheckBoxIndeterminate *widget.Icon
	ArrowLeft             *widget.Icon
	ArrowRight            *widget.Icon
	ArrowUp               *widget.Icon
	ArrowDown             *widget.Icon
}

func GetIcons() *Icons {
	iconsOnce.Do(func() {
		icns = new(Icons)
		icns.ContentAdd, _ = widget.NewIcon(icons.ContentAdd)
		icns.ActionDelete, _ = widget.NewIcon(icons.ActionDelete)
		icns.CheckBoxBlank, _ = widget.NewIcon(icons.ToggleCheckBoxOutlineBlank)
		icns.CheckBoxIndeterminate, _ = widget.NewIcon(icons.ToggleIndeterminateCheckBox)
		icns.ArrowLeft, _ = widget.NewIcon(icons.HardwareKeyboardArrowLeft)
		icns.ArrowRight, _ = widget.NewIcon(icons.HardwareKeyboardArrowRight)
		icns.ArrowUp, _ = widget.NewIcon(icons.HardwareKeyboardArrowUp)
		icns.ArrowDown, _ = widget.NewIcon(icons.HardwareKeyboardArrowDown)
	})
	return icns
}

package razgio

import (
	"gioui.org/layout"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Page interface {
	GetName() string
	Select()
	Layout(gtx C) D
}

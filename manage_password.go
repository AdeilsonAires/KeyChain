//go:generate fyne bundle -o bundled.go assets

package main

import (
	"fyne.io/fyne/v2"

	"KeyChain/layout"
)

func CreateManagePasswordWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("KeyChain")
	w.Resize(fyne.NewSize(400, 300))

	hw := &layout.LayoutWindow{ 
		App: a,
		Win: w,
	}

	hw.SetMainContainer("managePassword") 
	w.SetContent(hw.MainContainer)
	return w
}

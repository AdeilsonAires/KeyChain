//go:generate fyne bundle -o bundled.go assets

package main

import (
	"fyne.io/fyne/v2"

	"KeyChain/layout"
)

func CreateHomeWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("KeyChain")
	hw := &layout.LayoutWindow{ 
		App: a,
		Win: w,
	}

	hw.SetMainContainer("home")
	w.SetContent(hw.MainContainer)
	return w
}

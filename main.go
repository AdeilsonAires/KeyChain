//go:generate fyne bundle -o bundled.go assets
package main

import (
	database "KeyChain/dataBase"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const defaultWin = "home"


func CreateWindow(a fyne.App, winType string) fyne.Window {
	if winType == "home" {
		return CreateHomeWindow(a)
	} else if winType == "managePassword" {
		return CreateManagePasswordWindow(a)
	}
	return a.NewWindow("Janela Padrão")
}

func main() {
	database.CreateTable()
	a := app.NewWithID("KeyChain")
	a.SetIcon(resourceKeychainiconPng)
	win := CreateWindow(a, defaultWin)
	win.CenterOnScreen()
	win.ShowAndRun()
}

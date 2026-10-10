package main

/*
#include <gtk/gtk.h>
*/
import "C"
import (
	"log/slog"
)

//export HandleKeyPress
func HandleKeyPress(keyval uint, state uint) bool {
	switch currentMode { // call to mode.go
	case 0:		// ModeNormal
		switch keyval {
		case 105:	// i
			currentMode = ModeInsert
			appWindow.SetStatus("[INSERT] ")
			slog.Info("mode switched to INSERT")
		case 118:	// v
			currentMode = ModeVisible
			appWindow.SetStatus("[VISIBLE] ")
			slog.Info("mode switched to VISIBLE")
		case 106: 	// j
			slog.Info("detected j")
			if appWindow != nil {
				appWindow.Scroll(1) // scroll to down
			}
		case 107:	// k
			slog.Info("detected k")
			if appWindow != nil {
				appWindow.Scroll(-1) // scroll to up
			}
		case 111:	// o
			currentMode = ModeCommand
			appWindow.SetStatus("[COMMAND] ")
			slog.Info("detected o")
			//OpenUrl()
		}
		return true
	default:	// ModeInsert / ModeVisible
		if keyval == 0xff1b {		// escape
			currentMode = ModeNormal
			appWindow.SetStatus("[NORMAL] ")
			slog.Info("mode switched to NORMAL")
		}
		return false
	}
	return false
}

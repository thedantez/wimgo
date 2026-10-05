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
			slog.Info("mode switched to INSERT")
			return true
		case 118:	// v
			currentMode = ModeVisible
			slog.Info("mode switched to VISIBLE")
			return true
		case 106: 	// j
			slog.Info("detected j")
			if appWindow != nil {
				appWindow.Scroll(1) // scroll to down
			}
			return true
		case 107:	// k
			slog.Info("detected k")
			if appWindow != nil {
				appWindow.Scroll(-1) // scroll to up
			}
			return true
		}
	default:	// ModeInsert / ModeVisible
		if keyval == 0xff1b {		// escape
			currentMode = ModeNormal
			slog.Info("mode switched to NORMAL")
		}
		return false
	}
	return false
}

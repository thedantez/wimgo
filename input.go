package main

/*
#include <gtk/gtk.h>
*/
import "C"
import (
	"log/slog"
	//"os"
)

//export HandleKeyPress
func HandleKeyPress(keyval C.guint, state C.guint) C.gboolean {
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
			appWindow.OpenCommand(":open ")
			appWindow.SetStatus("[COMMAND] ")
			slog.Info("detected o")
		case 58:	// :
			currentMode = ModeCommand
			appWindow.OpenCommand(":")
			appWindow.SetStatus("[COMMAND] ")
			slog.Info("detected :")
		//case 113:	// q
		//	slog.Info("exitign...")
		//	//os.Exit(0)		// hard exit
		//	C.gtk_main_quit()	// soft exit
		case 114: 	// r
			appWindow.Reload()
			slog.Info("detected r")
		case 82:	// R
			appWindow.ReloadBypassCache()
			slog.Info("detected R")
		}
		return C.TRUE
	default:	// ModeInsert/ModeVisible/ModeCommand
		if keyval == 0xff1b {		// escape
			if currentMode == ModeCommand {
				appWindow.CloseCommand()
			}
			currentMode = ModeNormal
			appWindow.SetStatus("[NORMAL] ")
			slog.Info("mode switched to NORMAL")
		}
		return C.FALSE
	}
	return C.FALSE
}

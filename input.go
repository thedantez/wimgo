package main

/*
#include <gtk/gtk.h>
*/
import "C"
import (
	"time"
	"strconv"
	"log/slog"
	//"os"
)

var (
	keyBuffer 	string
	chordTimer  *time.Timer
	statusDirty bool
)

const chordTimeout = 1200 * time.Millisecond	// 1.2 seconds to finish combination

func resetKeyBuffer() {
	keyBuffer = ""
	if chordTimer != nil {
		chordTimer.Stop()
		chordTimer = nil
	}
	if statusDirty && appWindow != nil {
		statusDirty = false
		appWindow.SetStatus("[NORMAL] ")
	}
}

func splitCountAndCmd(buf string) (int, string) {
	idx := 0
	for idx < len(buf) && buf[idx] >= '0' && buf[idx] <= '9' {
		idx++
	}
	if idx == 0 {
		return 1, buf
	}
	count, err := strconv.Atoi(buf[:idx])
	if err != nil || count <= 0 {
		count = 1
	}
	return count, buf[idx:]
}

func processNormalChord(buf string) (done bool, pending bool) {
	count, cmd := splitCountAndCmd(buf)

	if cmd == "" || cmd == "g" {
		return false, true
	}

	switch cmd {
	case "j":
		appWindow.Scroll(0, count)
	case "k":
		appWindow.Scroll(0, -count)
	case "l":
		appWindow.Scroll(count, 0)
	case "h":
		appWindow.Scroll(-count, 0)
	case "gg":
		appWindow.ScrollToTop()
	case "G", "gG":
		appWindow.ScrollToBottom()
	case "r":
		appWindow.Reload()
	case "R":
		appWindow.ReloadBypassCache()
	case "H":
		appWindow.GoBack()
	case "L":
		appWindow.GoForward()
	default:
		slog.Warn("unknown key chord", "chord", buf)
	}
	return true, false
}

//export HandleKeyPress
func HandleKeyPress(keyval C.guint, state C.guint) C.gboolean {
	switch currentMode { // call to mode.go
		case 0:		// ModeNormal
		if keyval >= 0xffe1 && keyval <= 0xffee {	// ignore SHIFT_*, Contrlo_*
			return C.FALSE
		}
		if	keyval == 0xff1b {	// esc
			resetKeyBuffer()
			return C.TRUE
		}
		if keyBuffer == "" {
			switch keyval {
				case 105:	// i
					currentMode = ModeInsert
					appWindow.SetStatus("[INSERT] ")
					slog.Info("mode switched to INSERT")
					return C.TRUE
				case 118:	// v
					currentMode = ModeVisible
					appWindow.SetStatus("[VISIBLE] ")
					slog.Info("mode switched to VISIBLE")
					return C.TRUE
				//case 106: 	// j
				//slog.Info("detected j")
				//if appWindow != nil {
				//	appWindow.Scroll(0, 1) // scroll to down
				//}
				//case 107:	// k
				//slog.Info("detected k")
				//if appWindow != nil {
				//	appWindow.Scroll(0, -1) // scroll to up
				//}
				//case 108:	// l
				//slog.Info("detected l")
				//if appWindow != nil {
				//	appWindow.Scroll(1, 0)
				//}
				//case 104:	// h
				//slog.Info("detected h")
				//if appWindow != nil {
				//	appWindow.Scroll(-1, 0)
				//}
				case 111:	// o
					currentMode = ModeCommand
					appWindow.OpenCommand(":open ")
					appWindow.SetStatus("[COMMAND] ")
					slog.Info("detected o")
					return C.TRUE
				case 58:	// :
					currentMode = ModeCommand
					appWindow.OpenCommand(":")
					appWindow.SetStatus("[COMMAND] ")
					slog.Info("detected :")
					return C.TRUE
				//case 113:	// q
				//	slog.Info("exitign...")
				//	//os.Exit(0)		// hard exit
				//	C.gtk_main_quit()	// soft exit
				//case 114: 	// r
				//appWindow.Reload()
				//slog.Info("detected r")
				//case 82:	// R
				//appWindow.ReloadBypassCache()
				//slog.Info("detected R")
				//case 72:	// H
				//appWindow.GoBack()
				//slog.Info("detected H")
				//case 76:	// L
				//appWindow.GoForward()
				//slog.Info("detected L")
			}
		}
		if keyval >= 32 && keyval <= 126 {
			keyBuffer += string(rune(keyval))
			if chordTimer != nil {
				chordTimer.Stop()
			}

			done, pending := processNormalChord(keyBuffer)
			if done {
				resetKeyBuffer()
			} else if pending {
				statusDirty = true
				appWindow.SetStatus("[NORMAL] " + keyBuffer + " ")
				chordTimer = time.AfterFunc(chordTimeout, func() {
					resetKeyBuffer()
				})
			}
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

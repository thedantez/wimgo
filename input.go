package main

/*
#include <gtk/gtk.h>
*/
import "C"
import (
	"time"
	"strconv"
	"log/slog"
)

var (
	keyBuffer 	string
	chordTimer  *time.Timer
	statusDirty bool
)

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

	if cmd == "" || cmd == "g" || cmd == "z" {
		return false, true
	}

	action, exists := appConfig.KeyBindings[cmd]
	if !exists {
		slog.Warn("unknown or unbound key chord", "chord", buf)
		return true, false
	}

	switch action {
	case "scroll_down":
		appWindow.Scroll(0, count)
	case "scroll_up":
		appWindow.Scroll(0, -count)
	case "scroll_right":
		appWindow.Scroll(count, 0)
	case "scroll_left":
		appWindow.Scroll(-count, 0)
	case "scroll_top":
		appWindow.ScrollToTop()
	case "scroll_bottom":
		appWindow.ScrollToBottom()
	case "hints":
		StartHints()
	case "soft_reload":
		appWindow.Reload()
	case "hard_reload":
		appWindow.ReloadBypassCache()
	case "history_back":
		appWindow.GoBack()
	case "history_forward":
		appWindow.GoForward()
	case "zoom_in":
		appWindow.ZoomIn()
	case "zoom_out":
		appWindow.ZoomOut()
	case "reset_zoom":
		appWindow.ResetZoom()
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
			}
			resetKeyBuffer()
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
				chordTimer = time.AfterFunc(time.Duration(appConfig.ChordTimeout) * time.Millisecond, func() {
					resetKeyBuffer()
				})
			}
		}
		return C.TRUE
	case 4:		// ModeHints
		if keyval == 0xff1b {		// escape
			ExitHints()
			slog.Info("mode switched to NORMAL")
		}
		if keyval >= 32 && keyval <= 126 {
			ProcessHintKey(rune(keyval))
			return C.TRUE
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

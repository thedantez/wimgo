package main

/*
#include <gtk/gtk.h>
*/
import "C"

import(
	_ "embed"
	"fmt"
	"log/slog"
)

//go:embed hints.js
var hintsJS string

var (
	hintsInjected bool
	hintBuffer	  string
)

//export ResetHintsState
func ResetHintsState() {
	hintsInjected = false
	if currentMode == ModeHints {
		ExitHints()
	}
}

//export FinishHintAction
func FinishHintAction(targetMode C.int) {
	hintBuffer = ""
	if int(targetMode) == 1 {
		currentMode = ModeInsert
		if appWindow != nil {
			appWindow.SetStatus("[INSERT] ")
		}
		slog.Info("hint targeted input, mode switched to INSERT")
	} else {
		currentMode = ModeNormal
		if appWindow != nil {
			appWindow.SetStatus("[NORMAL] ")
		}
		slog.Info("hint completed, mode switched to NORMAL")
	}
}

func StartHints() {
	if appWindow == nil {
		return
	}

	hintBuffer = ""

	if !hintsInjected {
		appWindow.RunJS(hintsJS)
		hintsInjected = true
	}

	appWindow.RunJS("window.__wimgoHints && window.__wimgoHints.show();")
	currentMode = ModeHints
	appWindow.SetStatus("[HINTS] ")
	slog.Info("mode switched to HINTS")
}

func ProcessHintKey(char rune) {
	if appWindow == nil {
		return
	}

	hintBuffer += string(char)
	appWindow.SetStatus("[HINTS] " + hintBuffer + " ")

	jsCmd := fmt.Sprintf(`
	(function() {
		if (!window.__wimgoHints) return;
		var res = window.__wimgoHints.filter('%s');
		if (res === 'click:insert') {
			document.title = '__wimgo_hint_insert__';
		} else if (res === 'click:normal' || res === 'nomatch') {
			document.title = '__wimgo_hint_done__';
		}
	})();
	`, hintBuffer)
	appWindow.RunJS(jsCmd)
}

func ExitHints() {
	if appWindow != nil {
		appWindow.RunJS("window.__wimgoHints && window.__wimgoHints.clear();")
		appWindow.SetStatus("[NORMAL] ")
	}
	hintBuffer = ""
	currentMode = ModeNormal
	slog.Info("mode switched to NORMAL")
}

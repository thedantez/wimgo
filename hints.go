package main

/*
#include <gtk/gtk.h>
*/
import "C"

import (
	_ "embed"
	"fmt"
	"log/slog"
)

//go:embed hints.js
var hintsJS string

var hintBuffer string

//export ResetHintsState
func ResetHintsState() {
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
		slog.Info("hint targeted input: switched to INSERT")
	} else {
		currentMode = ModeNormal
		if appWindow != nil {
			appWindow.SetStatus("[NORMAL] ")
		}
		slog.Info("hint completed: switched to NORMAL")
	}
}

func StartHints() {
	if appWindow == nil {
		return
	}

	hintBuffer = ""
	currentMode = ModeHints
	appWindow.SetStatus("[HINTS] ")
	slog.Info("mode switched to HINTS")

	// Один атомарный вызов: инициализирует __wimgoHints и сразу внутри вызывает .show()
	appWindow.RunJS(hintsJS)
}

func ProcessHintKey(char rune) {
	if appWindow == nil {
		return
	}

	hintBuffer += string(char)
	appWindow.SetStatus("[HINTS] " + hintBuffer + " ")

	jsCmd := fmt.Sprintf(`
	(function() {
		if (!window.__wimgoHints) {
			document.title = '__wimgo_hint_done__';
			return;
		}
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

package main

/*
#include <gtk/gtk.h>
*/
import "C"
import (
	"strings"
	"log/slog"
	//"os"
)

func HandleCmdEntry(raw string) {	// takes raw from cmdEntry from window.go & parse it
	raw = strings.TrimSpace(raw)
	if len(raw) == 0 { 		// if command not exists
		if appWindow != nil {
			appWindow.CloseCommand()
			currentMode = ModeNormal
			appWindow.SetStatus("[NORMAL] ")
		}
		return
	}

	if strings.HasPrefix(raw, ":") {		// if command exists
		cleanCmd := strings.TrimPrefix(raw, ":")
		//split by :command & flags & parametr (what to open, for example)
		parts := strings.Fields(cleanCmd)	// except :open -t duckduckgo.com
		if len(parts) == 0 { //	commands not found
			return
		}

		cmd := parts[0]
		args := parts[1:]
		ExecuteCommand(cmd, args)
	} else {
		OpenUrl(raw)
	}

	if appWindow != nil {
		appWindow.CloseCommand()
		currentMode = ModeNormal
		appWindow.SetStatus("[NORMAL] ")
	}
}

func ExecuteCommand(cmd string, args []string) {
	slog.Info("handle command tooks command", "cmd", cmd, "args", args)
	switch cmd {
	case "q", "quit":
		slog.Info("executing quit command")
		//os.Exit(0)	// hard exit
		C.gtk_main_quit() 	// soft exit
	case "open", "o":
		if len(args) == 0 {
			slog.Warn("open command requires a url")
			return
		}
		target := strings.Join(args, " ")
		OpenUrl(target)
	case "reload", "r":
		if appWindow != nil {
			appWindow.Reload()
			slog.Info("reloading site...")
		}
	case "reload!", "r!":
		if appWindow != nil {
			appWindow.ReloadBypassCache()
			slog.Info("reloading site w/ bypass cache...")
		}
	case "b", "back":
		if appWindow != nil {
			appWindow.GoBack()
			slog.Info("go back history")
		}
	case "f", "forward":
		if appWindow != nil {
			appWindow.GoForward()
			slog.Info("go forward history")
		}
	default:
		slog.Warn("unknown command", "cmd", cmd)
	}
}

func OpenUrl(rawURL string) {
	slog.Info("opening url", "rawURL", rawURL)
	target := strings.TrimSpace(rawURL)
	if len(target) == 0 {
		return
	}

	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		if strings.Contains(target, " ") {
			target = "https://google.com/search?q=" + target
		} else {
			if strings.Contains(target, ".") && !strings.HasPrefix(target, ".") && !strings.HasSuffix(target, ".") {
				target = "https://" + target
			} else {
				target = "https://google.com/search?q=" + target
			}
		}
	}

	slog.Info("navigatin to URL", "url", target)
	if appWindow != nil {
		appWindow.LoadURL(target)
	}
}

package main

/*
#include <gtk/gtk.h>
*/
import "C"
import (
	"strings"
	"log/slog"
	"strconv"
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
	case "o", "open":
		if len(args) == 0 {
			slog.Warn("open command requires a url")
			return
		}
		target := strings.Join(args, " ")
		OpenUrl(target)
	case "r", "reload":
		if appWindow != nil {
			appWindow.Reload()
			slog.Info("reloading site...")
		}
	case "r!", "reload!":
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
	case "tabnew", "t":
		target := appConfig.HomePage
		if len(args) > 0 {
			target = strings.Join(args, " ")
		}
		if appWindow != nil {
			appWindow.NewTab(target)
		}
	case "tabclose", "tc":
		if appWindow != nil {
			appWindow.CloseTab(appWindow.activeTab)
		}
	case "tabnext", "tn":
		if appWindow != nil {
			appWindow.NextTab()
		}
	case "tabprev", "tp":
		if appWindow != nil {
			appWindow.PrevTab()
		}
	case "set":
		if len(args) < 2 {
			slog.Warn("command set requires key & value")
			return
		}
		switch args[0] {
			case "scrollStep", "scrollstep", "scroll_step":
				if val, err := strconv.Atoi(args[1]); err == nil && val > 0 {
					appConfig.ScrollStep = val
					slog.Info("set config parameter", "param", "scrollStep", "val", val)
				}
			case "homePage", "homepage", "home_page":
				appConfig.HomePage = args[1]
				slog.Info("set config parameter", "param", "homePage", "val", args[1])
			case "searchEngine", "searchengine", "search_engine":
				appConfig.SearchEngine = args[1]
				slog.Info("set config parameter", "param", "searchEngine", "val", args[1])
			case "chordTimeout", "chordtimeout", "chord_timeout":
				if val, err := strconv.Atoi(args[1]); err == nil && val > 0 {
					appConfig.ChordTimeout = val
					slog.Info("set config parameter", "param", "chordTimeout", "val", val)
				} else {
					slog.Warn("invalid value for chordTimeout", "arg", args[1])
				}
			case "webgl", "webGL", "webGl":
				appConfig.WebGL = (args[1] == "true" || args[1] == "1" || args[1] == "on")
				if appWindow != nil {
					appWindow.SetWebGL(appConfig.WebGL)
				}
			default:
				slog.Info("unknown args for command set ", "value", args[0], "args", args[1])
		}
	case "bind":
		if len(args) < 2 {
			slog.Warn("bind command requires key & action (usage: bind <key> <action>)")
			return
		}
		key := args[0]
		action := args[1]

		if appConfig.KeyBindings == nil {
			appConfig.KeyBindings = make(map[string]string)
		}
		appConfig.KeyBindings[key] = action
		slog.Info("key bound successfully", "key", key, "action", action)
	case "unbind":
		if len(args) < 1 {
			slog.Warn("unbind command requires a key (usage: :unbind <key>)")
			return
		}
		key := args[0]
		delete(appConfig.KeyBindings, key)
		slog.Info("key unbound (disabled)", "key", key)
	case "bind-default", "reset-default":
		key := ""
		if len(args) > 0 {
			key = args[0]
		}
		ResetBinding(key)
	case "config-save":
		SaveConfig()
		slog.Info("config saved")
	case "config-load":
		LoadConfig()
		slog.Info("config loaded")
	case "stop":
		appWindow.StopLoading()
		slog.Info("stop loading")
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
			target = appConfig.SearchEngine + target
		} else {
			if strings.Contains(target, ".") && !strings.HasPrefix(target, ".") && !strings.HasSuffix(target, ".") {
				target = "https://" + target
			} else {
				target = appConfig.SearchEngine + target
			}
		}
	}

	slog.Info("navigatin to URL", "url", target)
	if appWindow != nil {
		appWindow.LoadURL(target)
	}
}

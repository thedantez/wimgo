package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"log/slog"
)

type Config struct {
	HomePage 		string
	SearchEngine	string
	ScrollStep		int
	ChordTimeout	int
	ZoomStep		float64
	WebGL			bool
	KeyBindings map[string]string
}

var appConfig = Config{		// default
	HomePage: 		"https://duckduckgo.com",
	SearchEngine:	"https://duckduckgo.com/?q=",
	ScrollStep:		100,
	ChordTimeout:	1200,
	ZoomStep:		0.25,
	WebGL:			false,
	KeyBindings: DefaultKeyBindings(),
}

func DefaultKeyBindings() map[string]string{
	return	map[string]string{
		"j":	"scroll_down",
		"k":	"scroll_up",
		"gg":	"scroll_top",
		"G":	"scroll_bottom",
		//"gG":	"scroll_bottom",
		"r":	"soft_reload",
		"R":	"hard_reload",
		"H":	"history_back",
		"L":	"history_forward",
		"J":	"tab_next",
		"K":	"tab_prev",
		"d":	"tab_close",
		"f":	"hints",
		"+":	"zoom_in",
		//"zi":	"zoom_in",
		"-":	"zoom_out",
		//"zo":	"zoom_out",
		"=":	"reset_zoom",
		//"zz":	"reset_zoom",
	}
}

func ResetBinding(key string) {
	defaults := DefaultKeyBindings()

	if key == "" || key == "all" {
		appConfig.KeyBindings = defaults
		slog.Info("all keybindings reset to defaults")
		return
	}

	if defaultAction, exists := defaults[key]; exists {
		appConfig.KeyBindings[key] = defaultAction
		slog.Info("keybinding reset to default", "key", key, "action", defaultAction)
	} else {
		delete(appConfig.KeyBindings, key)
		slog.Info("custom keybinding removed", "key", key)
	}
}

func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(configDir, "wimgo")
	return filepath.Join(configPath, "config"), nil
}

func LoadConfig() {
	configPath, err := getConfigPath()
	if err != nil {
		slog.Warn("cannot determine user config dit, using defaults", "err", err)
	}

	file, err := os.Open(configPath)
	if os.IsNotExist(err) {
		slog.Info("config not found, generating default config", "path", configPath)
		SaveConfig()
		return
	} else if err != nil {
		slog.Warn("failed to open config", "err", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 {
			ExecuteCommand(parts[0], parts[1:])
		}
	}
	slog.Info("config loaded successfully")
}

func SaveConfig() {
	configPath, err := getConfigPath()
	if err != nil {
		slog.Error("cannot get config path for saving config", "err", err)
		return
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		slog.Warn("failed to create config directory", "err", err)
		return
	}
	content := fmt.Sprintf(`
# wimgo configuration file
# this file is valuated line-by-line using wimgo commands

set homePage %s
set searchEngine %s
set scrollStep %d
set chordTimeout %d
set webgl %t`,
						appConfig.HomePage,
						appConfig.SearchEngine,
						appConfig.ScrollStep,
						appConfig.ChordTimeout,
						appConfig.WebGL)
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		slog.Error("failed to write config file", "err", err)
		return
	}
	slog.Info("config saved to ", "path", configPath)
}

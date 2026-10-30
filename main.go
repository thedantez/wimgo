package main

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.1
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>
*/
import "C"
import (
	"log/slog"
	"flag"
	"os"
	"path/filepath"
)

func main() {
	slog.Info("starting wimgo...")

	noCookies := flag.Bool("no-cookies", false, "run browser in ephemeral in-memory mode w/out saving cookies")
	flag.Parse()
	LoadConfig()

	var dataDir, cacheDir string
	if !*noCookies {
		home, _ := os.UserHomeDir()
		dataDir = filepath.Join(home, ".local", "share", "wimgo")
		cacheDir = filepath.Join(home, ".cache", "wimgo")
		_ = os.MkdirAll(dataDir, 0755)
		_ = os.MkdirAll(cacheDir, 0755)
	}

	win :=	NewWindow(appConfig.HomePage, dataDir, cacheDir, *noCookies)
	win.Show()

	slog.Info("window opened, entering into main cycle GTK events")
}

package main

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.1
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>
*/
import "C"
import (
	"log/slog"
)

func main() {
	slog.Info("starting wimgo...")

	LoadConfig()
	win :=	NewWindow(appConfig.HomePage)
	win.Show()

	slog.Info("window opened, entering into main cycle GTK events")
}

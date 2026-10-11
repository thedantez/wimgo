package main

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.1
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>
*/
import "C"
import (
	"log/slog"
	"os"
)

//type BrowserWindow struct {
//	window 	*C.GtkWidget
//	webView *C.WebKitWebView
//}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("starting wimgo...")

	win :=	NewWindow("https://duckduckgo.com")
	win.Show()

	logger.Info("window opened, entering into main cycle GTK events")
}

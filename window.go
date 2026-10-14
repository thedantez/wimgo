package main

// CGO & C-world

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.1
#include "callbacks.h"
*/
import "C" 	// import C-world

import ( 	// import Go's libs
	"log/slog"
	"fmt"
	"os"
	"runtime"
	"unsafe"	// for c-world injections
)

type BrowserWindow struct {
	window 		*C.GtkWidget
	webView 	*C.WebKitWebView
	statusLabel *C.GtkWidget
	cmdEntry	*C.GtkWidget
	infoLabel	*C.GtkWidget
}


// func for hide cmdEntry
func (b *BrowserWindow) Show() {
	C.gtk_widget_show_all(b.window)
	C.gtk_widget_hide(b.cmdEntry)
	C.gtk_main()
}


// func for reload sites
func (b *BrowserWindow) Reload() {		// like soft reload (use saved data from cache if it wasn't changed)
	C.webkit_web_view_reload(b.webView)
}

func (b *BrowserWindow) ReloadBypassCache() {	// like hard reload (full reload site)
	C.webkit_web_view_reload_bypass_cache(b.webView)
}


// func(-s) to scrolling
func (b *BrowserWindow) Scroll(deltaY int) {		// need to add scroll by Ox: create new func or modify this func
	jsScript := C.CString(fmt.Sprintf("window.scrollBy(0, %d)", (deltaY * 100)))
	defer C.free(unsafe.Pointer(jsScript))
	C.webkit_web_view_run_javascript(
		b.webView,
		jsScript,
		nil, nil, nil)
}

// func for modify statusLabel by normal/insert/command/visual
func (b *BrowserWindow) SetStatus(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.gtk_label_set_text((*C.GtkLabel)(unsafe.Pointer(b.statusLabel)), cText)
}


// funcs for ModeCommand by cmdEntry
func (b *BrowserWindow) OpenCommand(prefix string) {
	cPrefix := C.CString(prefix)
	defer C.free(unsafe.Pointer(cPrefix))

	// setting cmdEntry
	C.gtk_widget_show(b.cmdEntry)	// make possible to see cmdEntry
	C.gtk_entry_set_text(
		(*C.GtkEntry)(unsafe.Pointer(b.cmdEntry)),
		cPrefix,
	)
	C.gtk_widget_grab_focus(b.cmdEntry)
	C.gtk_editable_set_position(
		(*C.GtkEditable)(unsafe.Pointer(b.cmdEntry)),
		-1,
	)
}

func (b *BrowserWindow) CloseCommand() {
	cEmpty := C.CString("")
	defer C.free(unsafe.Pointer(cEmpty))
	C.gtk_entry_set_text(
		(*C.GtkEntry)(unsafe.Pointer(b.cmdEntry)),
		cEmpty,
	)
	C.gtk_widget_hide(b.cmdEntry)
	C.gtk_widget_grab_focus(
		(*C.GtkWidget)(unsafe.Pointer(b.webView)),
	)
}

//export HandleCmdActivate
func HandleCmdActivate(cText *C.char) {
	raw := C.GoString(cText)
	HandleCmdEntry(raw)
}

func (b *BrowserWindow) LoadURL(url string) {
	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))
	C.webkit_web_view_load_uri(b.webView, cURL)
}

// funcs for infoLabel
func (b *BrowserWindow) SetInfo(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.gtk_label_set_text(
		(*C.GtkLabel)(unsafe.Pointer(b.infoLabel)),
		cText,
	)
}

//export UpdateLoadProgress
func UpdateLoadProgress(progress float64) {
	if appWindow == nil {
		return
	}

	percent := int(progress * 100)
	if percent >= 100 {
		cURI := C.webkit_web_view_get_uri(appWindow.webView)
		if cURI != nil {
			appWindow.SetInfo(C.GoString(cURI))
		}
	} else {
		appWindow.SetInfo(fmt.Sprintf("[%d%%]", percent))
	}
}

// func for start GTK
func init() {
	runtime.LockOSThread()
	os.Setenv("WEBKIT_DISABLE_DMABUF_RENDERER", "1")
}

var appWindow *BrowserWindow	// global pointer for changes in window & ui

func NewWindow(url string) *BrowserWindow {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))		// creating stuctured logs
	logger.Info("starting creating window...")

	// main window GTK
	C.gtk_init(nil, nil)	// initializing GTK
	C.enable_dark_theme() 	// setting dark theme

	window := C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)	// creating window
	C.gtk_window_set_default_size((*C.GtkWindow)(unsafe.Pointer(window)), 1280, 720)

	keyPress := C.CString("key-press-event")	// key press handler
	defer C.free(unsafe.Pointer(keyPress))
	C.g_signal_connect_data(	// connect w/ C-world
		C.gpointer(window),
		(*C.gchar)(unsafe.Pointer(keyPress)),
		C.GCallback(C.on_key_press),
		nil, nil, 0,
	)

	cTitle := C.CString("wimgo browser")	// title of window
	defer C.free(unsafe.Pointer(cTitle))
	C.gtk_window_set_title((*C.GtkWindow)(unsafe.Pointer(window)), cTitle)

	// engine: WebKit
	webView := C.webkit_web_view_new()	// window for view web

	cSignalFailed := C.CString("load-failed")		// if web can't load
	defer C.free(unsafe.Pointer(cSignalFailed))
	C.g_signal_connect_data(	// connect w/ C-world
		C.gpointer(webView),
		(*C.gchar)(unsafe.Pointer(cSignalFailed)),
		C.GCallback(C.on_load_failed),
		nil, nil, 0,
	)

	cSignalChanged := C.CString("load-changed")		// if web loaded/changed
	defer C.free(unsafe.Pointer(cSignalChanged))
	C.g_signal_connect_data(	// connect w/ C-world
		C.gpointer(webView),
		(*C.gchar)(unsafe.Pointer(cSignalChanged)),
		C.GCallback(C.on_load_changed),
		nil, nil, 0,
	)

	cSignalProgress := C.CString("notify::estimated-load-progress")		// editing status of loading page in infoLabel
	defer C.free(unsafe.Pointer(cSignalProgress))
	C.g_signal_connect_data(
		C.gpointer(webView),
		(*C.gchar)(unsafe.Pointer(cSignalProgress)),
		C.GCallback(C.on_notify_progress),
		nil, nil, 0,
	)

	cSignalDestroy := C.CString("destroy")	// destroy when window was closed
	defer C.free(unsafe.Pointer(cSignalDestroy))
	C.g_signal_connect_data(	// connect w/ C-world
		C.gpointer(window),
		(*C.gchar)(unsafe.Pointer(cSignalDestroy)),
		C.GCallback(C.gtk_main_quit),
		nil, nil, 0,
	)

	// creating bottom (status, command) bar, infolabel
	cStatus := C.CString("[NORMAL] ")
	defer C.free(unsafe.Pointer(cStatus))
	statusLabel := C.gtk_label_new(cStatus)
	C.gtk_widget_set_halign((*C.GtkWidget)(unsafe.Pointer(statusLabel)), C.GTK_ALIGN_START)

	cmdEntry := C.gtk_entry_new()
	C.gtk_widget_hide(cmdEntry)		// hide cmdEntry
	C.gtk_entry_set_has_frame(		// edit cmdEntry's frame
		(*C.GtkEntry)(unsafe.Pointer(cmdEntry)),
		0,
	)

	cSignalActivate := C.CString("activate")	// handler cmdEntry for entering commands
	defer C.free(unsafe.Pointer(cSignalActivate))
	C.g_signal_connect_data(
		C.gpointer(cmdEntry),
		(*C.gchar)(unsafe.Pointer(cSignalActivate)),
		C.GCallback(C.on_cmd_activate),
		nil, nil, 0,
	)

	cURLText := C.CString(url)
	defer C.free(unsafe.Pointer(cURLText))
	infoLabel := C.gtk_label_new(cURLText)
	C.gtk_label_set_ellipsize(
		(*C.GtkLabel)(unsafe.Pointer(infoLabel)),
		C.PANGO_ELLIPSIZE_END,
	)
	C.gtk_label_set_max_width_chars(
		(*C.GtkLabel)(unsafe.Pointer(infoLabel)),
		65,
	)
	C.gtk_label_set_single_line_mode(	// single mode for label
		(*C.GtkLabel)(unsafe.Pointer(infoLabel)),
		C.TRUE,
	)

	// style by CSS
	cssData := `
		window, box {
			background-color: #000000;
		}
		entry, label {
			border: none;
			box-shadow: none;
			background: transparent;
			font-family: monospace;
			font-size: 13px;
			min-height: 0px;
			margin: 0px;
			padding: 2px 4px;
		}
		#bottom-bar {
			color: #ffffff;
			background-color: #000000;
			border-top: 1px solid #313244;
			padding: 2px 4px;
		}
		#bottom-bar label {
			color: #ffffff;
			font-family: monospace;
			font-size: 13px;
			padding: 0 4px;
		}
		#bottom-bar entry {
			background: transparent;
			color: #ffffff;
			font-family: monospace;
			font-size: 13px;
			border: none;
			box-shadow: none;
			padding: 0;
			margin: 0;
		}
	`
	cCSSData := C.CString(cssData)
	defer C.free(unsafe.Pointer(cCSSData))

	cssProvider := C.gtk_css_provider_new()
	C.gtk_css_provider_load_from_data(
		cssProvider,
		cCSSData,
		-1,
		nil,
	)

	screen := C.gdk_screen_get_default()
	C.gtk_style_context_add_provider_for_screen(
		screen,
		(*C.GtkStyleProvider)(unsafe.Pointer(cssProvider)),
		C.GTK_STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
	C.g_object_unref(C.gpointer(cssProvider))

	// style: enabling dark theme & dark mode
	C.apply_webview_dark_mode((*C.WebKitWebView)(unsafe.Pointer(webView)))	// enabling dark scheme

	// packing
	mainbox := C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0)		// main box
	bottombox := C.gtk_box_new(C.GTK_ORIENTATION_HORIZONTAL, 6)		// bottom box w/ status,url,
	cBottomBarID := C.CString("bottom-bar")
	defer C.free(unsafe.Pointer(cBottomBarID))
	C.gtk_widget_set_name(
		(*C.GtkWidget)(unsafe.Pointer(bottombox)),
		cBottomBarID,
	)

	C.gtk_container_add(
		(*C.GtkContainer)(unsafe.Pointer(window)),	// packing main box to window
		(*C.GtkWidget)(unsafe.Pointer(mainbox)),
	)

	C.gtk_box_pack_start(	// packing webView window to main box
		(*C.GtkBox)(unsafe.Pointer(mainbox)),
		(*C.GtkWidget)(unsafe.Pointer(webView)),
		1,	// expand = true
		1,	// fill = true
		0,	// padding = 0
	)

	C.gtk_box_pack_start(	// packing statusLabel to bottom box
		(*C.GtkBox)(unsafe.Pointer(bottombox)),
		(*C.GtkWidget)(unsafe.Pointer(statusLabel)),
		0,	// expand = false
		0,	// fill = false
		0,	// padding = 0
	)

	C.gtk_box_pack_start(	// packing cmdEntry input area to bottom box
		(*C.GtkBox)(unsafe.Pointer(bottombox)),
		(*C.GtkWidget)(unsafe.Pointer(cmdEntry)),
		1,	// expand = false
		1,	// fill = false
		0,	// padding = 0
	)

	C.gtk_box_pack_end(		// packing infoLabel to bottom box
		(*C.GtkBox)(unsafe.Pointer(bottombox)),
		(*C.GtkWidget)(unsafe.Pointer(infoLabel)),
		0,	// expand = false
		0,	// fill = false
		4,	// padding = 0
	)

	C.gtk_box_pack_end(		// packing bottombox to main box
		(*C.GtkBox)(unsafe.Pointer(mainbox)),
		(*C.GtkWidget)(unsafe.Pointer(bottombox)),
		0,	// expand = false
		0,	// fill = false
		0,	// padding = 0
	)

	// loading url
	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))

	C.webkit_web_view_load_uri(
		(*C.WebKitWebView)(unsafe.Pointer(webView)),
		cURL,
	)

	// creating appWindow global variable & pack all
	appWindow = &BrowserWindow{		// modify appWindow variable for return & save changes in global level
		window: 		window,
		webView: 		(*C.WebKitWebView)(unsafe.Pointer(webView)),
		statusLabel: 	(*C.GtkWidget)(unsafe.Pointer(statusLabel)),
		cmdEntry:		(*C.GtkWidget)(unsafe.Pointer(cmdEntry)),
		infoLabel: 		(*C.GtkWidget)(unsafe.Pointer(infoLabel)),
	}
	return appWindow
}

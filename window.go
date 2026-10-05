package main

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.1
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>
gboolean on_load_failed(WebKitWebView *web_view,
								WebKitLoadEvent load_event,
								gchar *failing_uri,
								GError *error,
								gpointer user_data) {
	g_print("[WebKit Error] URI: %s | Error: %s\n", failing_uri, error->message);
	return FALSE;
}

void on_load_changed(WebKitWebView *web_view,
							 WebKitLoadEvent load_event,
						 	 gpointer user_data) {
	if (load_event == WEBKIT_LOAD_FINISHED) {
		g_print("[WebKit] Page load finished: %s\n", webkit_web_view_get_uri(web_view));
	}
}

extern int HandleKeyPress(unsigned int keyval, unsigned int state);
gboolean on_key_press(GtkWidget *widget, GdkEventKey *event, gpointer user_data) {
	g_print("[KeyPress] key pressed %u : %u\n", event->keyval, event->state);
	return HandleKeyPress(event->keyval, event->state);
}
*/
import "C"
import (
	"log/slog"
	"fmt"
	"os"
	"runtime"
	"unsafe"
)

type BrowserWindow struct {
	window 		*C.GtkWidget
	webView 	*C.WebKitWebView
	statusLabel *C.GtkWidget
}

func (b *BrowserWindow) Show() {
	C.gtk_widget_show_all(b.window)
	C.gtk_main()
}

func (b *BrowserWindow) Scroll(deltaY int) {
	jsScript := C.CString(fmt.Sprintf("window.scrollBy(0, %d)", (deltaY * 100)))
	defer C.free(unsafe.Pointer(jsScript))
	C.webkit_web_view_run_javascript(
		b.webView,
		jsScript,
		nil, nil, nil)
}

func (b *BrowserWindow) SetStatus(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.gtk_label_set_text((*C.GtkLabel)(unsafe.Pointer(b.statusLabel)), cText)
}

func init() {
	runtime.LockOSThread()
}

var appWindow *BrowserWindow

func NewWindow(url string) *BrowserWindow {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("starting creating window...")

	C.gtk_init(nil, nil)

	window := C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
	C.gtk_window_set_default_size((*C.GtkWindow)(unsafe.Pointer(window)), 1280, 720)
	keyPress := C.CString("key-press-event")
	defer C.free(unsafe.Pointer(keyPress))
	C.g_signal_connect_data(
		C.gpointer(window),
		(*C.gchar)(unsafe.Pointer(keyPress)),
		C.GCallback(C.on_key_press),
		nil, nil, 0,
	)

	box := C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0)
	C.gtk_container_add(
		(*C.GtkContainer)(unsafe.Pointer(window)),
		(*C.GtkWidget)(unsafe.Pointer(box)),
	)

	cTitle := C.CString("wimgo browser")
	defer C.free(unsafe.Pointer(cTitle))
	C.gtk_window_set_title((*C.GtkWindow)(unsafe.Pointer(window)), cTitle)

	webView := C.webkit_web_view_new()

	cSignalFailed := C.CString("load-failed")
	defer C.free(unsafe.Pointer(cSignalFailed))
	C.g_signal_connect_data(
		C.gpointer(webView),
		(*C.gchar)(unsafe.Pointer(cSignalFailed)),
		C.GCallback(C.on_load_failed),
		nil, nil, 0,
	)

	cSignalChanged := C.CString("load-changed")
	defer C.free(unsafe.Pointer(cSignalChanged))
	C.g_signal_connect_data(
		C.gpointer(webView),
		(*C.gchar)(unsafe.Pointer(cSignalChanged)),
		C.GCallback(C.on_load_changed),
		nil, nil, 0,
	)

	C.gtk_box_pack_start(
		(*C.GtkBox)(unsafe.Pointer(box)),
		(*C.GtkWidget)(unsafe.Pointer(webView)),
		1,	// expand = true
		1,	// fill = true
		0,	// padding = 0
	)

	cStatus := C.CString("[NORMAL]")
	defer C.free(unsafe.Pointer(cStatus))
	statusLabel := C.gtk_label_new(cStatus)
	C.gtk_widget_set_halign((*C.GtkWidget)(unsafe.Pointer(statusLabel)), C.GTK_ALIGN_END)
	C.gtk_box_pack_start(
		(*C.GtkBox)(unsafe.Pointer(box)),
		(*C.GtkWidget)(unsafe.Pointer(statusLabel)),
		0,	// expand = false
		0,	// fill = false
		0,	// padding = 0
	)

	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))

	C.webkit_web_view_load_uri(
		(*C.WebKitWebView)(unsafe.Pointer(webView)),
		cURL,
	)

	cSignalDestroy := C.CString("destroy")
	defer C.free(unsafe.Pointer(cSignalDestroy))
	C.g_signal_connect_data(
		C.gpointer(window),
		(*C.gchar)(unsafe.Pointer(cSignalDestroy)),
		C.GCallback(C.gtk_main_quit),
		nil, nil, 0,
	)

	appWindow = &BrowserWindow{
		window: window,
		webView: (*C.WebKitWebView)(unsafe.Pointer(webView)),
		statusLabel: (*C.GtkWidget)(unsafe.Pointer(statusLabel)),
	}
	return appWindow
}

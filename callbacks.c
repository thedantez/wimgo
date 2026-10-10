#include "callbacks.h"
#include "_cgo_export.h"

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

gboolean on_key_press(GtkWidget *widget,
                      GdkEventKey *event,
                      gpointer user_data) {
	g_print("[KeyPress] key pressed %u : %u\n", event->keyval, event->state);
	return HandleKeyPress(event->keyval, event->state);
}

void on_notify_progress(WebKitWebView *web_view,
                        GParamSpec *pspec,
                        gpointer user_data) {
	double progress = webkit_web_view_get_estimated_load_progress(web_view);
	UpdateLoadProgress(progress);
}

void enable_dark_theme(void) {
    GtkSettings *settings = gtk_settings_get_default();
    if (settings != NULL) {
        g_object_set(settings, "gtk-application-prefer-dark-theme", TRUE, NULL);
    }
}

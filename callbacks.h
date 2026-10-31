#ifndef CALLBACKS_H
#define CALLBACKS_H

#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

WebKitWebView* create_custom_webview(const char *data_dir,
                                     const char *cache_dir,
                                     gboolean ephemeral);

void set_webview_webgl(WebKitWebView *web_view, gboolean flag);

gboolean on_load_failed(WebKitWebView *web_view,
								WebKitLoadEvent load_event,
								gchar *failing_uri,
								GError *error,
								gpointer user_data);

void on_load_changed(WebKitWebView *web_view,
							 WebKitLoadEvent load_event,
						 	 gpointer user_data);

gboolean on_key_press(GtkWidget *widget, GdkEventKey *event, gpointer user_data);

void on_notify_progress(WebKitWebView *web_view, GParamSpec *pspec, gpointer user_data);

void on_window_resized(GtkWidget *widget, GParamSpec *pspec, gpointer user_data);
void on_cmd_activate(GtkEntry *entry, gpointer user_data);

// dark theme, scheme, mode
void enable_dark_theme(void);

void apply_webview_dark_mode(WebKitWebView *web_view);

void on_notify_title(WebKitWebView *web_view, GParamSpec *pspec, gpointer user_data);

void inject_js(WebKitWebView *web_view, const char *js_code);
#endif

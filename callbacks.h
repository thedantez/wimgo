#ifndef CALLBACKS_H
#define CALLBACKS_H

#include <gtk/gtk.h>
#include <webkit2/webkit2.h>
#include <stdlib.h>

void init_global_context(const char *data_dir,
                         const char *cache_dir,
                         gboolean ephemeral);

WebKitWebView* create_custom_webview(void);

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

gboolean on_window_resized(GtkWidget *widget, GdkEventConfigure *event, gpointer user_data);

void on_cmd_activate(GtkEntry *entry, gpointer user_data);

// dark theme, scheme, mode
void enable_dark_theme(void);

void apply_webview_dark_mode(WebKitWebView *web_view);

void on_notify_title(WebKitWebView *web_view, GParamSpec *pspec, gpointer user_data);

void inject_js(WebKitWebView *web_view, const char *js_code);

void print_page(WebKitWebView *web_view);

extern void HandleDownloadRequest(char* url, char* suggested_filename);

static gboolean cancel_download_idle(gpointer data);

void on_download_started(WebKitWebContext *context,
                         WebKitDownload *download,
                         gpointer user_data);

void schedule_label_update(GtkWidget *widget, const char *text, int show);
#endif

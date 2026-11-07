#include "callbacks.h"
#include "_cgo_export.h"

static WebKitWebContext *global_context = NULL;

void init_global_context(const char *data_dir,
                         const char *cache_dir,
                         gboolean ephemeral) {
    if (global_context != NULL) {
        return;
    }

    WebKitWebsiteDataManager *manager;
    if (ephemeral) {
        manager = webkit_website_data_manager_new_ephemeral();
    } else {
        manager = webkit_website_data_manager_new(
                "base-data-directory", data_dir,
                "base-cache-directory", cache_dir,
                NULL
                );
        WebKitCookieManager *cookie_mgr = webkit_website_data_manager_get_cookie_manager(manager);
        gchar *cookie_file = g_build_filename(data_dir, "cookie.db", NULL);
        webkit_cookie_manager_set_persistent_storage(
                cookie_mgr,
                cookie_file,
                WEBKIT_COOKIE_PERSISTENT_STORAGE_SQLITE
                );
        webkit_cookie_manager_set_accept_policy(
                cookie_mgr,
                WEBKIT_COOKIE_POLICY_ACCEPT_ALWAYS
                );
        g_free(cookie_file);
    }

    global_context = webkit_web_context_new_with_website_data_manager(manager);
    g_object_unref(manager);

    g_signal_connect(global_context,    // for download dialog
                    "download-started",
                    G_CALLBACK(on_download_started),
                    NULL);
}

WebKitWebView* create_custom_webview() {
    if (global_context == NULL) {
        return WEBKIT_WEB_VIEW(webkit_web_view_new());
    }
    return WEBKIT_WEB_VIEW(webkit_web_view_new_with_context(global_context));
}

// code for download
struct label_update_data {
	GtkLabel *label;
	GtkWidget *widget;
	char *text;
	int show;
};

static gboolean update_label_idle(gpointer user_data) {
	struct label_update_data *data = (struct label_update_data *)user_data;
	if (data->show) {
		gtk_widget_show(data->widget);
		gtk_label_set_text(data->label, data->text);
	} else {
		gtk_widget_hide(data->widget);
	}
	free(data->text);
	free(data);
	return FALSE;
}

void schedule_label_update(GtkWidget *widget, const char *text, int show) {
	struct label_update_data *data = malloc(sizeof(struct label_update_data));
	data->widget = widget;
	data->label = GTK_LABEL(widget);
	data->text = strdup(text);
	data->show = show;
	g_idle_add(update_label_idle, data);
}

gboolean on_load_failed(WebKitWebView *web_view,
                        WebKitLoadEvent load_event,
                        gchar *failing_uri,
                        GError *error,
                        gpointer user_data) {
	g_print("[WebKit Error] URI: %s | Error: %s\n", failing_uri, error->message);
    UpdateTabState(web_view, NULL, 0, 1);
	return FALSE;
}

void on_load_changed(WebKitWebView *web_view,
                     WebKitLoadEvent load_event,
                     gpointer user_data) {
	if (load_event == WEBKIT_LOAD_FINISHED) {
		g_print("[WebKit] Page load finished: %s\n", webkit_web_view_get_uri(web_view));
        UpdateTabState(web_view, (char*)webkit_web_view_get_title(web_view), 1, 0);
	}
}

gboolean on_key_press(GtkWidget *widget,
                      GdkEventKey *event,
                      gpointer user_data) {
	//g_print("[KeyPress] key pressed %u : %u\n", event->keyval, event->state);
	return HandleKeyPress(event->keyval, event->state);
}

void on_notify_progress(WebKitWebView *web_view,
                        GParamSpec *pspec,
                        gpointer user_data) {
	double progress = webkit_web_view_get_estimated_load_progress(web_view);
	UpdateLoadProgress(progress);
}

gboolean on_window_resized(GtkWidget *widget,
                       GdkEventConfigure *event,
                       gpointer user_data) {
    TriggerUpdateTabBar();
    return FALSE;
}

void on_notify_title(WebKitWebView *web_view, GParamSpec *pspec, gpointer user_data) {
    const gchar *title = webkit_web_view_get_title(web_view);
    if (title == NULL) return;
    if (title && !g_str_has_prefix(title, "__wimgo_hint_")) {
        UpdateTabState(web_view, (char*)title, 0, 0);
    }

    if (g_strcmp0(title, "__wimgo_hint_done__") == 0) {
        FinishHintAction(0); // 0 = ModeNormal
        return;
    }
    if (g_strcmp0(title, "__wimgo_hint_insert__") == 0) {
        FinishHintAction(1); // 1 = ModeInsert
        return;
    }
}

void on_cmd_activate(GtkEntry *entry,
                     gpointer user_data) {
    const gchar *text = gtk_entry_get_text(entry);
    if (text != NULL) {
        HandleCmdActivate((char *)text);
    }
}

void enable_dark_theme(void) {
    GtkSettings *settings = gtk_settings_get_default();
    if (settings != NULL) {
        g_object_set(settings, "gtk-application-prefer-dark-theme", TRUE, NULL);
    }
}

void set_webview_webgl(WebKitWebView *web_view, gboolean flag) {
    if (!web_view) return;
    WebKitSettings *settings = webkit_web_view_get_settings(web_view);
    if (settings != NULL) {
        webkit_settings_set_enable_webgl(settings, flag);
    }
}

void apply_webview_dark_mode(WebKitWebView *web_view) {     // 0 =no_preference, 1=light, 2=dark
    WebKitSettings *settings = webkit_web_view_get_settings(web_view);
    WebKitWebContext *context = webkit_web_view_get_context(web_view);
    webkit_web_context_set_cache_model(context, WEBKIT_CACHE_MODEL_DOCUMENT_VIEWER);
    if (settings != NULL) {
        if (g_object_class_find_property(G_OBJECT_GET_CLASS(settings), "color-scheme")) {
            g_object_set(G_OBJECT(settings), "color-scheme", 2, NULL);
        }
        webkit_settings_set_enable_media_stream(settings, TRUE);
        webkit_settings_set_enable_mediasource(settings, TRUE);
        webkit_settings_set_hardware_acceleration_policy(settings, WEBKIT_HARDWARE_ACCELERATION_POLICY_ON_DEMAND);
    }
}

void inject_js(WebKitWebView *web_view, const char *js_code) {
    WebKitUserContentManager *manager = webkit_web_view_get_user_content_manager(web_view);
    if (!manager) return;

    WebKitUserScript *script = webkit_user_script_new(
            js_code,
            WEBKIT_USER_CONTENT_INJECT_TOP_FRAME,   // or WEBKIT_USER_CONTENT_INJECT_ALL_FRAMES
            WEBKIT_USER_SCRIPT_INJECT_AT_DOCUMENT_END,
            NULL,
            NULL
            );
    webkit_user_content_manager_add_script(manager, script);
    webkit_user_script_unref(script);
}

void print_page(WebKitWebView *web_view) {
    WebKitPrintOperation *op = webkit_print_operation_new(web_view);
    webkit_print_operation_run_dialog(op, NULL);
    g_object_unref(op);
}

extern void HandleDownloadRequest(char* url, char* suggested_filename);

static gboolean cancel_download_idle(gpointer data) {
    WebKitDownload *dl = WEBKIT_DOWNLOAD(data);
    webkit_download_cancel(dl);
    g_object_unref(dl);
    return FALSE;
}

void on_download_started(WebKitWebContext *context,
                         WebKitDownload *download,
                         gpointer user_data) {
    WebKitURIRequest *req= webkit_download_get_request(download);
    const gchar *uri = webkit_uri_request_get_uri(req);

    if (uri != NULL) {
        HandleDownloadRequest((char*)uri, (char*)"");
    }

    g_object_ref(download);
    g_idle_add(cancel_download_idle, download);
}

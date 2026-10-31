package main

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.1
#include <webkit2/webkit2.h>
*/
import "C"

type TabState int

const (
	TabLoading TabState = iota
	TabLoaded
	TabFailed
)

type BrowserTab struct {
	id			int
	webView 	*C.WebKitWebView
	title		string
	uri			string
	state		TabState
	progress	int
}

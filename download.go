package main

/*
#cgo pkg-config: gtk+-3.0
#include "callbacks.h"
*/
import "C"
import (
	"fmt"
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type DownloadState int

const (
	DlDownloading DownloadState = iota
	DlFinished
	DlFailed
)

type DownloadItem struct {
	Name 		string
	Path		string
	Progress	int
	State		DownloadState
}

var (
	pendingDownloadURL		string
	pendingDownloadPath		string
	dlMu					sync.Mutex
	downloads 				[]*DownloadItem
)

//export HandleDownloadRequest
func HandleDownloadRequest(cURL, cFilename *C.char) {
	rawURL := C.GoString(cURL)
	suggestedName := C.GoString(cFilename)

	if suggestedName == "" {
		if strings.HasPrefix(rawURL, "data:image/") {
			ext := "jpg"
			if strings.HasPrefix(rawURL, "data:image/png") {
				ext = "png"
			} else if strings.HasPrefix(rawURL, "data:image/webp") {
				ext = "webp"
			}
			suggestedName = "image." + ext
		} else {
			parts := strings.Split(strings.Split(rawURL, "?")[0], "/")
			if len(parts) > 0 {
				suggestedName = parts[len(parts) - 1]
			}
			if suggestedName == "" {
				suggestedName = "downloaded_file"
			}
		}
	}

	pendingDownloadURL = rawURL

	home, _ := os.UserHomeDir()
	defaultPath := filepath.Join(home, "Downloads", suggestedName)

	if appWindow != nil {
		currentMode = ModeDownloadPath
		appWindow.SetStatus("[DOWNLOAD to] ")
		appWindow.OpenCommand(defaultPath)
	}
}

func cancelDownload() {
	pendingDownloadURL = ""
	pendingDownloadPath = ""
	if appWindow != nil {
		appWindow.CloseCommand()
		currentMode = ModeNormal
		appWindow.SetStatus("[NORMAL] ")
	}
}

func checkDownloadPath() {
	if info, err := os.Stat(pendingDownloadPath); err == nil && info.IsDir() {
		cancelDownload()
		return
	}

	dir := filepath.Dir(pendingDownloadPath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		appWindow.SetStatus(fmt.Sprintf("[CREATE DIR %s? y/n] ", dir))
		appWindow.OpenCommand("")
		currentMode = ModeDownloadConfirmDir
		return
	}

	checkDownloadFileOverwrite()
}

func checkDownloadFileOverwrite() {
	if _, err := os.Stat(pendingDownloadPath); err == nil {
		appWindow.SetStatus("[OVERWRITE FILE? y/n] ")
		appWindow.OpenCommand("")
		currentMode = ModeDownloadConfirmOverwrite
		return
	}
	startDownload()
}

func startDownload() {
	url := pendingDownloadURL
	dest := pendingDownloadPath

	cancelDownload()

	if info, err := os.Stat(dest); err == nil && info.IsDir() {
		return
	}

	item := &DownloadItem{
		Name:		filepath.Base(dest),
		Path:		dest,
		Progress:	0,
		State:		DlDownloading,
	}

	dlMu.Lock()
	downloads = append(downloads, item)
	dlMu.Unlock()
	UpdateDownloadBar()

	slog.Info("starting download", "url", url, "dest", dest)

	go func() {
		err := performDownload(url, dest, item)
		dlMu.Lock()
		if err != nil {
			slog.Warn("downlaod failed", "err", err)
			item.State = DlFailed
		} else {
			slog.Info("downlaod finished", "dest", dest)
			item.State = DlFinished
		}
		dlMu.Unlock()
		UpdateDownloadBar()
	}()
}

func performDownload(url, dest string, item *DownloadItem) error {
	if strings.HasPrefix(url, "data:") {
		commaIdx := strings.IndexByte(url, ',')
		if commaIdx == -1 {
			return fmt.Errorf("invalid data URI")
		}
		meta := url[:commaIdx]
		rawData := url[commaIdx+1:]
		var bytesData []byte
		var err error
		if strings.Contains(meta, ";base64") {
			bytesData, err = base64.StdEncoding.DecodeString(rawData)
			if err != nil {
				return err
			}
		} else {
			bytesData = []byte(rawData)
		}
		return os.WriteFile(dest, bytesData, 0644)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	total := resp.ContentLength
	var downloaded int64
	buf := make([]byte, 32*1024)
	lastPercent := -1

	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
			downloaded += int64(n)
			if total > 0 {
				percent := int((downloaded * 100) / total)
				if percent != lastPercent {
					lastPercent = percent
					dlMu.Lock()
					item.Progress = percent
					dlMu.Unlock()
					UpdateDownloadBar()
				}
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return readErr
		}
	}
	return nil
}

func UpdateDownloadBar() {
	if appWindow != nil {
		appWindow.UpdateDownloadBar()
	}
}

func ClearDownloads() {
	dlMu.Lock()
	for _, d := range downloads {
		if d.State == DlFailed || d.State == DlDownloading {
			if d.Path != "" {
				os.Remove(d.Path)
			}
			cancelDownload()
		}
	}
	downloads = nil
	dlMu.Unlock()

	UpdateDownloadBar()
}

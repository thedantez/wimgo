package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"html"
)

type HistoryEntry struct {
	Timestamp	time.Time
	URL			string
	Title		string
}

type HistoryManager struct {
	mu			sync.Mutex
	filePath	string
	lastURL		string
}

var globalHistory *HistoryManager

func InitHistory() {
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Warn("cannot determine user home for history", "err", err)
		return
	}

	histDir := filepath.Join(home, ".local", "share", "wimgo")
	if err := os.MkdirAll(histDir, 0755); err != nil {
		slog.Warn("failed to create history directory", "err", err)
	}

	globalHistory = &HistoryManager{
		filePath: filepath.Join(histDir, "history"),
	}
}

func (h *HistoryManager) AddEntry(rawURL, title string) {
	if h == nil || rawURL == "" || rawURL == "about:blank" || strings.HasPrefix(rawURL, "wimgo://") {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if rawURL == h.lastURL {
		return
	}
	h.lastURL = rawURL

	if title == "" {
		title = rawURL
	}

	title = strings.ReplaceAll(title, "\t", " ")
	title = strings.ReplaceAll(title, "\n", " ")

	now := time.Now()
	line := fmt.Sprintf("%s\t%s\t%s\n", now.Format(time.RFC3339), rawURL, title)

	f, err := os.OpenFile(h.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Warn("failed to append history line", "err", err)
	}
	defer f.Close()

	if _, err := f.WriteString(line); err != nil {
		slog.Warn("failed to append history line", "err", err)
	}
}

func (h *HistoryManager) GetRecent(limit int) []HistoryEntry {
	if h == nil {
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	f, err := os.Open(h.filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var entries []HistoryEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")
		if len(parts) >= 2 {
			t, _ := time.Parse(time.RFC3339, parts[0])
			title := parts[1]
			if len(parts) >= 3 {
				title = parts[2]
			}
			entries = append(entries, HistoryEntry{
				Timestamp: t,
				URL:	   parts[1],
				Title:	   title,
			})
		}
	}

	if len(entries) > limit {
		return entries[len(entries)-limit:]
	}
	return entries
}

func (h *HistoryManager) RenderHistoryRows() string {
	entries := h.GetRecent(100)

	var rows strings.Builder
	for i := len(entries) - 1; i >= 0; i-- {
		item := entries[i]
		escapedTitle := html.EscapeString(item.Title)
		escapedURL := html.EscapeString(item.URL)
		timeStr := item.Timestamp.Format("02.01 15:04")

		rows.WriteString(fmt.Sprintf(`
			<div class="entry">
				<span class="time">%s</span>
				<a class="title" href="%s">%s</a>
				<span class="url">%s</span>
			</div>`, timeStr, escapedURL, escapedTitle, escapedURL))
	}
	return rows.String()
}

func (h *HistoryManager) RenderHistoryHTML() string {
	entries := h.GetRecent(100)
	rowsHTML := h.RenderHistoryRows()

	return fmt.Sprintf(`<!DOCTYPE html>
	<html>
	<head>
	<meta charset="utf-8">
	<title>wimgo://history</title>
	<style>
	body {
		background-color: #000000;
		color: #e0e0e0;
		font-family: monospace;
		margin: 0;
		padding: 24px 32px;
	}
	h1 {
		font-size: 14px;
		color: #ffffff;
		border-bottom: 1px solid #333333;
		padding-bottom: 8px;
		margin-bottom: 16px;
		letter-spacing: 1px;
		text-transform: uppercase;
	}
	.entry {
		display: flex;
		align-items: baseline;
		padding: 5px 0;
		border-bottom: 1px solid #141414;
		font-size: 13px;
	}
	.time {
		color: #666666;
		width: 100px;
		flex-shrink: 0;
	}
	.title {
		color: #ffffff;
		text-decoration: none;
		margin-right: 16px;
		max-width: 480px;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}
	.title:hover {
		color: #aaaaaa;
		text-decoration: underline;
	}
	.url {
		color: #555555;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}
	</style>
	</head>
	<body>
	<h1>History (%d entries)</h1>
	<div id="history-list">
	%s
	</div>
	</body>
	</html>`, len(entries), rowsHTML)
}

package ui

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

const (
	historyFileName = "history.json"
	maxHistoryItems = 100
)

// historyEntry represents a single history item for persistence
type historyEntry struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Body    string            `json:"body,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

// hash returns a unique hash for the entire request
func (e historyEntry) hash() string {
	// Create a deterministic representation
	// Sort header keys for consistent ordering
	var headerKeys []string
	for k := range e.Headers {
		headerKeys = append(headerKeys, k)
	}
	sort.Strings(headerKeys)

	// Build a canonical string representation
	h := sha256.New()
	h.Write([]byte(e.Method))
	h.Write([]byte(e.URL))
	h.Write([]byte(e.Body))
	for _, k := range headerKeys {
		h.Write([]byte(k))
		h.Write([]byte(e.Headers[k]))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// getDataDir returns the path to ~/.getboy, creating it if needed
func getDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".getboy")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

// loadHistory reads history from disk
func loadHistory() ([]historyEntry, error) {
	dir, err := getDataDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, historyFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No history yet
		}
		return nil, err
	}

	var entries []historyEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// saveHistory writes history to disk
func saveHistory(entries []historyEntry) error {
	dir, err := getDataDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, historyFileName)
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// addToHistory adds a new entry to history, avoiding duplicates
// Uniqueness is based on the full request hash (method, URL, body, headers)
// Returns the updated history list
func addToHistory(entries []historyEntry, entry historyEntry) []historyEntry {
	entryHash := entry.hash()

	// Check for duplicate based on full request hash
	for i, e := range entries {
		if e.hash() == entryHash {
			// Move existing entry to front (most recent)
			entries = append(entries[:i], entries[i+1:]...)
			break
		}
	}

	// Prepend new entry
	entries = append([]historyEntry{entry}, entries...)

	// Trim to max size
	if len(entries) > maxHistoryItems {
		entries = entries[:maxHistoryItems]
	}

	return entries
}

// historyToItems converts history entries to list items
func historyToItems(entries []historyEntry) []reqItem {
	items := make([]reqItem, len(entries))
	for i, e := range entries {
		// Create a short title from the URL
		title := e.Method + " " + truncateURL(e.URL, 30)
		items[i] = reqItem{
			title:   title,
			desc:    e.URL,
			method:  e.Method,
			url:     e.URL,
			body:    e.Body,
			headers: e.Headers,
		}
	}
	return items
}

// truncateURL shortens a URL for display
func truncateURL(u string, maxLen int) string {
	if len(u) <= maxLen {
		return u
	}
	return u[:maxLen-3] + "..."
}

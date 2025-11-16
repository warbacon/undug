package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Bang struct {
	T string `json:"t"`
	U string `json:"u"`
}

type Bangs []Bang

var (
	bangs    Bangs
	bangsMap map[string]*Bang
)

func fetchBangs() error {
	fmt.Println("Fetching bangs from DuckDuckGo...")

	resp, err := http.Get("https://duckduckgo.com/bang.js")
	if err != nil {
		return fmt.Errorf("error fetching bangs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	err = json.Unmarshal(body, &bangs)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	bangsMap = make(map[string]*Bang, len(bangs))
	for i := range bangs {
		key := strings.ToLower(bangs[i].T)
		bangsMap[key] = &bangs[i]
	}

	fmt.Printf("Loaded %d bangs\n", len(bangs))

	return nil
}

func findBang(trigger string) *Bang {
	trigger = strings.TrimPrefix(trigger, "!")
	trigger = strings.ToLower(trigger)
	return bangsMap[trigger]
}

func logRedirect(start time.Time, action, query, url string) {
	elapsed := time.Since(start)
	fmt.Printf("[%s] %s (%.2fms)\n",
		time.Now().Format("15:04:05"),
		action,
		float64(elapsed.Microseconds())/1000.0)
	fmt.Printf("  Query: %s\n", query)
	fmt.Printf("  URL:   %s\n", url)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "No query provided", http.StatusBadRequest)
		return
	}

	words := strings.Fields(query)
	var redirectURL string

	for i, word := range words {
		if !strings.HasPrefix(word, "!") {
			continue
		}

		if bang := findBang(word); bang != nil {
			searchParts := append(words[:i], words[i+1:]...)
			searchTerm := strings.Join(searchParts, " ")
			redirectURL = strings.ReplaceAll(bang.U, "{{{s}}}", searchTerm)
			logRedirect(start, "Bang redirect", query, redirectURL)
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}
	}

	redirectURL = "https://google.com/search?q=" + query
	logRedirect(start, "Google fallback", query, redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func main() {
	err := fetchBangs()
	if err != nil {
		fmt.Println("Error loading bangs:", err)
		os.Exit(1)
	}

	http.HandleFunc("/", handleRequest)

	fmt.Println("Server listening on http://localhost:8765")
	err = http.ListenAndServe(":8765", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

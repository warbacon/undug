package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Bang struct {
	T string `json:"t"`
	U string `json:"u"`
}

type Bangs []Bang

var bangs Bangs

func loadBangs(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	err = json.Unmarshal(data, &bangs)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	fmt.Printf("Loaded %d bangs\n", len(bangs))
	return nil
}

func findBang(trigger string) *Bang {
	trigger = strings.TrimPrefix(trigger, "!")
	trigger = strings.ToLower(trigger)

	for _, bang := range bangs {
		if strings.ToLower(bang.T) == trigger {
			return &bang
		}
	}
	return nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "No query provided", http.StatusBadRequest)
		return
	}

	words := strings.Fields(query)

	for i, word := range words {
		if strings.HasPrefix(word, "!") {
			trigger := word

			searchParts := append(words[:i], words[i+1:]...)
			searchTerm := strings.Join(searchParts, " ")

			bang := findBang(trigger)
			if bang != nil {
				url := strings.ReplaceAll(bang.U, "{{{s}}}", searchTerm)
				http.Redirect(w, r, url, http.StatusFound)
				return
			}
		}
	}

	http.Redirect(w, r, "https://google.com/search?q="+query, http.StatusFound)
}

func main() {
	err := loadBangs("bangs.json")
	if err != nil {
		fmt.Println("Error loading bangs:", err)
		fmt.Println("Make sure bangs.json exists in the current directory")
		os.Exit(1)
	}

	http.HandleFunc("/", handleRequest)

	fmt.Println("Server listening on http://localhost:8765")

	err = http.ListenAndServe(":8765", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

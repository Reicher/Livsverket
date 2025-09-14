package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

// Set represents a simplified taxonomic set returned to the frontend
// It mirrors the structure requested by the acceptance criteria.
type Set struct {
	ID             string `json:"id"`
	ScientificName string `json:"scientificName"`
	Rank           string `json:"rank"`
	ParentID       string `json:"parentId"`
	ChildCount     int    `json:"childCount"`
}

// cacheMutex protects the childCache map.
var cacheMutex sync.RWMutex
var childCache = make(map[string][]Set)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/sets/root", handleRoot)
	mux.HandleFunc("/sets/", handleSets)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

// handleRoot serves the root sets from sets.json.
func handleRoot(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile(filepath.Join("..", "sets.json"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// handleSets dispatches requests under /sets/ .
func handleSets(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/sets/")
	if strings.HasSuffix(path, "/children") {
		id := strings.TrimSuffix(path, "/children")
		handleChildren(w, r, id)
		return
	}

	http.Error(w, "not found", http.StatusNotFound)
}

// handleChildren fetches children for a taxon either from cache or the
// Catalogue of Life API.
func handleChildren(w http.ResponseWriter, r *http.Request, id string) {
	sets, err := getChildren(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sets)
}

// getChildren returns the direct children of a COL taxon id. The result is
// cached in memory to avoid repeated network requests.
func getChildren(id string) ([]Set, error) {
	cacheMutex.RLock()
	if sets, ok := childCache[id]; ok {
		cacheMutex.RUnlock()
		return sets, nil
	}
	cacheMutex.RUnlock()

	// Fetch from Catalogue of Life API
	url := fmt.Sprintf("https://api.catalogueoflife.org/dataset/3LR/taxon/%s/children", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("COL API status: %s", resp.Status)
	}

	var colResp struct {
		Result []struct {
			ID   string `json:"id"`
			Name struct {
				ScientificName string `json:"scientificName"`
			} `json:"name"`
			Rank       string `json:"rank"`
			ParentID   string `json:"parentId"`
			ChildCount int    `json:"childCount"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&colResp); err != nil {
		return nil, err
	}

	sets := make([]Set, 0, len(colResp.Result))
	for _, r := range colResp.Result {
		sets = append(sets, Set{
			ID:             r.ID,
			ScientificName: r.Name.ScientificName,
			Rank:           r.Rank,
			ParentID:       r.ParentID,
			ChildCount:     r.ChildCount,
		})
	}

	cacheMutex.Lock()
	childCache[id] = sets
	cacheMutex.Unlock()

	return sets, nil
}

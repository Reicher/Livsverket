package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Set struct {
	ID             string `json:"id"`
	ScientificName string `json:"scientificName"`
	Rank           string `json:"rank"`
	ParentID       string `json:"parentId"`
	ChildCount     int    `json:"childCount"`
	Description    string `json:"description,omitempty"`
}

type Sighting struct {
	SetID     string    `json:"setId"`
	Timestamp time.Time `json:"timestamp"`
	Note      string    `json:"note,omitempty"`
}

type Server struct {
	sets      map[string]Set
	children  map[string][]Set
	idAliases map[string]string
	dataDir   string
	mu        sync.Mutex
}

const datasetKey = "3LR"
const colBaseURL = "https://api.catalogueoflife.org"

func (s *Server) resolveID(id string) (string, error) {
	if col, ok := s.idAliases[id]; ok {
		if col != "" {
			return col, nil
		}
		set := s.sets[id]
		cid, err := s.searchCOLID(set.ScientificName)
		if err != nil {
			return "", err
		}
		s.idAliases[id] = cid
		set.ID = cid
		s.sets[cid] = set
		return cid, nil
	}
	return id, nil
}

func (s *Server) searchCOLID(name string) (string, error) {
	u := fmt.Sprintf("%s/dataset/%s/nameusage/search?limit=1&q=%s", colBaseURL, datasetKey, url.QueryEscape(name))
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("COL search status: %s", resp.Status)
	}
	var res struct {
		Result []struct {
			ID string `json:"id"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	if len(res.Result) == 0 {
		return "", errors.New("not found")
	}
	return res.Result[0].ID, nil
}

func (s *Server) fetchSet(id string) (Set, error) {
	if set, ok := s.sets[id]; ok {
		return set, nil
	}
	colID, err := s.resolveID(id)
	if err != nil {
		return Set{}, err
	}
	if set, ok := s.sets[colID]; ok {
		return set, nil
	}
	u := fmt.Sprintf("%s/dataset/%s/taxon/%s", colBaseURL, datasetKey, colID)
	resp, err := http.Get(u)
	if err != nil {
		return Set{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Set{}, fmt.Errorf("COL taxon status: %s", resp.Status)
	}
	var col struct {
		ID         string `json:"id"`
		Rank       string `json:"rank"`
		ParentID   string `json:"parentId"`
		ChildCount int    `json:"childCount"`
		Name       struct {
			ScientificName string `json:"scientificName"`
		} `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&col); err != nil {
		return Set{}, err
	}
	set := Set{ID: col.ID, ScientificName: col.Name.ScientificName, Rank: col.Rank, ParentID: col.ParentID, ChildCount: col.ChildCount}
	s.sets[col.ID] = set
	return set, nil
}

func (s *Server) fetchChildren(id string) ([]Set, error) {
	colID, err := s.resolveID(id)
	if err != nil {
		return nil, err
	}
	if list, ok := s.children[colID]; ok {
		return list, nil
	}
	u := fmt.Sprintf("%s/dataset/%s/taxon/%s/children", colBaseURL, datasetKey, colID)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("COL children status: %s", resp.Status)
	}
	var res struct {
		Result []struct {
			ID         string `json:"id"`
			Rank       string `json:"rank"`
			ParentID   string `json:"parentId"`
			ChildCount int    `json:"childCount"`
			Name       struct {
				ScientificName string `json:"scientificName"`
			} `json:"name"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	list := make([]Set, len(res.Result))
	for i, c := range res.Result {
		child := Set{ID: c.ID, ScientificName: c.Name.ScientificName, Rank: c.Rank, ParentID: c.ParentID, ChildCount: c.ChildCount}
		s.sets[child.ID] = child
		list[i] = child
	}
	s.children[colID] = list
	if alias, ok := s.idAliases[id]; ok && alias == colID {
		s.children[id] = list
	}
	return list, nil
}

func loadSets(path string) (map[string]Set, map[string][]Set, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	var list []Set
	if err := json.NewDecoder(f).Decode(&list); err != nil {
		return nil, nil, err
	}
	sets := make(map[string]Set)
	children := make(map[string][]Set)
	for _, s := range list {
		sets[s.ID] = s
		parent := s.ParentID
		children[parent] = append(children[parent], s)
	}
	return sets, children, nil
}

func NewServer(dataDir string) (*Server, error) {
	sets, children, err := loadSets(filepath.Join(dataDir, "sets.json"))
	if err != nil {
		return nil, err
	}
	aliases := make(map[string]string)
	for id := range sets {
		aliases[id] = ""
	}
	return &Server{sets: sets, children: children, idAliases: aliases, dataDir: dataDir}, nil
}

func (s *Server) writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (s *Server) getSet(id string) (Set, bool) {
	set, ok := s.sets[id]
	return set, ok
}

func (s *Server) handleRootSets(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, s.children["root"])
}

func (s *Server) handleSet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	set, err := s.fetchSet(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	s.writeJSON(w, set)
}

func (s *Server) handleChildren(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	list, err := s.fetchChildren(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, list)
}

func (s *Server) handleBreadcrumbs(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var chain []Set
	for id != "" && id != "root" {
		set, err := s.fetchSet(id)
		if err != nil {
			break
		}
		chain = append([]Set{set}, chain...)
		id = set.ParentID
	}
	s.writeJSON(w, chain)
}

func (s *Server) readCollection() ([]string, error) {
	path := filepath.Join(s.dataDir, "collections.json")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var ids []string
	if err := json.NewDecoder(f).Decode(&ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func (s *Server) writeCollection(ids []string) error {
	path := filepath.Join(s.dataDir, "collections.json")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(ids)
}

func (s *Server) handleGetCollections(w http.ResponseWriter, r *http.Request) {
	ids, err := s.readCollection()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, ids)
}

func (s *Server) handlePostCollection(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SetID string `json:"setId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	ids, err := s.readCollection()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, id := range ids {
		if id == req.SetID {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	ids = append(ids, req.SetID)
	if err := s.writeCollection(ids); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleDeleteCollection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.mu.Lock()
	defer s.mu.Unlock()
	ids, err := s.readCollection()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var out []string
	for _, existing := range ids {
		if existing != id {
			out = append(out, existing)
		}
	}
	if err := s.writeCollection(out); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) readSightings() ([]Sighting, error) {
	path := filepath.Join(s.dataDir, "sightings.json")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var list []Sighting
	if err := json.NewDecoder(f).Decode(&list); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Server) writeSightings(list []Sighting) error {
	path := filepath.Join(s.dataDir, "sightings.json")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(list)
}

func (s *Server) handleGetSightings(w http.ResponseWriter, r *http.Request) {
	list, err := s.readSightings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, list)
}

func (s *Server) handlePostSighting(w http.ResponseWriter, r *http.Request) {
	var req Sighting
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	list, err := s.readSightings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	list = append(list, req)
	if err := s.writeSightings(list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func main() {
	srv, err := NewServer(".")
	if err != nil {
		log.Fatalf("failed to load data: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/sets/root", srv.handleRootSets)
	mux.HandleFunc("/sets/{id}", srv.handleSet)
	mux.HandleFunc("/sets/{id}/children", srv.handleChildren)
	mux.HandleFunc("/sets/{id}/breadcrumbs", srv.handleBreadcrumbs)

	mux.HandleFunc("/collections", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			srv.handleGetCollections(w, r)
		case http.MethodPost:
			srv.handlePostCollection(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/collections/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			srv.handleDeleteCollection(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/sightings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			srv.handleGetSightings(w, r)
		case http.MethodPost:
			srv.handlePostSighting(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := 8080
	fmt.Printf("Server running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

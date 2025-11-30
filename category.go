package main

import (
	"encoding/json"
	"net/http"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var cat Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	cat.ID = store.nextCatID
	store.nextCatID++
	store.categories[cat.ID] = &cat
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	cats := make([]*Category, 0, len(store.categories))
	for _, c := range store.categories {
		cats = append(cats, c)
	}
	store.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cats)
}

func GetCategoryByID(id int) (*Category, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	cat, exists := store.categories[id]
	return cat, exists
}

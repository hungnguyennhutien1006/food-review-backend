package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	w.WriteHeader(http.StatusCreated)
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

func GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	store.mu.RLock()
	cat, exists := store.categories[id]
	store.mu.RUnlock()

	if !exists {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cat)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	store.mu.Lock()
	cat, exists := store.categories[id]
	if !exists {
		store.mu.Unlock()
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	var updatedCat Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCat); err != nil {
		store.mu.Unlock()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedCat.ID = cat.ID
	store.categories[id] = &updatedCat
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCat)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	store.mu.Lock()
	_, exists := store.categories[id]
	if !exists {
		store.mu.Unlock()
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	delete(store.categories, id)
	store.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func GetCategoryByID(id int) (*Category, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	cat, exists := store.categories[id]
	return cat, exists
}

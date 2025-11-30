package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Restaurant struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	OwnerID   int       `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var rest Restaurant
	if err := json.NewDecoder(r.Body).Decode(&rest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, userExists := GetUserByID(rest.OwnerID)
	if !userExists {
		http.Error(w, "Owner (user) not found", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	rest.ID = store.nextRestID
	store.nextRestID++
	rest.CreatedAt = time.Now()
	store.restaurants[rest.ID] = &rest
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rest)
}

func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	rests := make([]*Restaurant, 0, len(store.restaurants))
	for _, r := range store.restaurants {
		rests = append(rests, r)
	}
	store.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rests)
}

func GetRestaurant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	store.mu.RLock()
	rest, exists := store.restaurants[id]
	store.mu.RUnlock()

	if !exists {
		http.Error(w, "Restaurant not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rest)
}

func UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	store.mu.Lock()
	rest, exists := store.restaurants[id]
	if !exists {
		store.mu.Unlock()
		http.Error(w, "Restaurant not found", http.StatusNotFound)
		return
	}

	var updatedRest Restaurant
	if err := json.NewDecoder(r.Body).Decode(&updatedRest); err != nil {
		store.mu.Unlock()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate owner exists
	_, userExists := GetUserByID(updatedRest.OwnerID)
	if !userExists {
		store.mu.Unlock()
		http.Error(w, "Owner (user) not found", http.StatusBadRequest)
		return
	}

	updatedRest.ID = rest.ID
	updatedRest.CreatedAt = rest.CreatedAt
	store.restaurants[id] = &updatedRest
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedRest)
}

func DeleteRestaurant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	store.mu.Lock()
	_, exists := store.restaurants[id]
	if !exists {
		store.mu.Unlock()
		http.Error(w, "Restaurant not found", http.StatusNotFound)
		return
	}

	delete(store.restaurants, id)
	store.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func GetRestaurantByID(id int) (*Restaurant, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	rest, exists := store.restaurants[id]
	return rest, exists
}
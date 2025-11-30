package main

import (
	"encoding/json"
	"net/http"
	"time"
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

func GetRestaurantByID(id int) (*Restaurant, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	rest, exists := store.restaurants[id]
	return rest, exists
}

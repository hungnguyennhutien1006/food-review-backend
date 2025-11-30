package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Food struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	CategoryID   int     `json:"category_id"`
	RestaurantID int     `json:"restaurant_id"`
}

func CreateFood(w http.ResponseWriter, r *http.Request) {
	var food Food
	if err := json.NewDecoder(r.Body).Decode(&food); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, catExists := GetCategoryByID(food.CategoryID)
	if !catExists {
		http.Error(w, "Category not found", http.StatusBadRequest)
		return
	}

	_, restExists := GetRestaurantByID(food.RestaurantID)
	if !restExists {
		http.Error(w, "Restaurant not found", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	food.ID = store.nextFoodID
	store.nextFoodID++
	store.foods[food.ID] = &food
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(food)
}

func GetFoods(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	foods := make([]*Food, 0, len(store.foods))
	for _, f := range store.foods {
		foods = append(foods, f)
	}
	store.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foods)
}

func GetFoodsByRestaurant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	restID, _ := strconv.Atoi(vars["restaurant_id"])

	store.mu.RLock()
	foods := make([]*Food, 0)
	for _, f := range store.foods {
		if f.RestaurantID == restID {
			foods = append(foods, f)
		}
	}
	store.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foods)
}

func GetFoodByID(id int) (*Food, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	food, exists := store.foods[id]
	return food, exists
}
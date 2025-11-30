package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	InitStorage()
	SeedData()

	r := mux.NewRouter()

	// User routes
	r.HandleFunc("/users", CreateUser).Methods("POST")
	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", GetUser).Methods("GET")

	// Category routes
	r.HandleFunc("/categories", CreateCategory).Methods("POST")
	r.HandleFunc("/categories", GetCategories).Methods("GET")

	// Restaurant routes
	r.HandleFunc("/restaurants", CreateRestaurant).Methods("POST")
	r.HandleFunc("/restaurants", GetRestaurants).Methods("GET")

	// Food routes
	r.HandleFunc("/foods", CreateFood).Methods("POST")
	r.HandleFunc("/foods", GetFoods).Methods("GET")
	r.HandleFunc("/restaurants/{restaurant_id}/foods", GetFoodsByRestaurant).Methods("GET")

	// Comment routes
	r.HandleFunc("/comments", CreateComment).Methods("POST")
	r.HandleFunc("/foods/{food_id}/comments", GetCommentsByFood).Methods("GET")
	r.HandleFunc("/users/{user_id}/comments", GetCommentsByUser).Methods("GET")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

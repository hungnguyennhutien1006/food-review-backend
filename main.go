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
	r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

	// Category routes
	r.HandleFunc("/categories", CreateCategory).Methods("POST")
	r.HandleFunc("/categories", GetCategories).Methods("GET")
	r.HandleFunc("/categories/{id}", GetCategory).Methods("GET")
	r.HandleFunc("/categories/{id}", UpdateCategory).Methods("PUT")
	r.HandleFunc("/categories/{id}", DeleteCategory).Methods("DELETE")

	// Restaurant routes
	r.HandleFunc("/restaurants", CreateRestaurant).Methods("POST")
	r.HandleFunc("/restaurants", GetRestaurants).Methods("GET")
	r.HandleFunc("/restaurants/{id}", GetRestaurant).Methods("GET")
	r.HandleFunc("/restaurants/{id}", UpdateRestaurant).Methods("PUT")
	r.HandleFunc("/restaurants/{id}", DeleteRestaurant).Methods("DELETE")

	// Food routes
	r.HandleFunc("/foods", CreateFood).Methods("POST")
	r.HandleFunc("/foods", GetFoods).Methods("GET")
	r.HandleFunc("/foods/{id}", GetFood).Methods("GET")
	r.HandleFunc("/foods/{id}", UpdateFood).Methods("PUT")
	r.HandleFunc("/foods/{id}", DeleteFood).Methods("DELETE")
	r.HandleFunc("/restaurants/{restaurant_id}/foods", GetFoodsByRestaurant).Methods("GET")

	// Comment routes
	r.HandleFunc("/comments", CreateComment).Methods("POST")
	r.HandleFunc("/comments/{id}", GetComment).Methods("GET")
	r.HandleFunc("/comments/{id}", UpdateComment).Methods("PUT")
	r.HandleFunc("/comments/{id}", DeleteComment).Methods("DELETE")
	r.HandleFunc("/foods/{food_id}/comments", GetCommentsByFood).Methods("GET")
	r.HandleFunc("/users/{user_id}/comments", GetCommentsByUser).Methods("GET")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
package main

import (
	"sync"
	"time"
)

type Storage struct {
	users       map[int]*User
	categories  map[int]*Category
	restaurants map[int]*Restaurant
	foods       map[int]*Food
	comments    map[int]*Comment
	mu          sync.RWMutex
	nextUserID  int
	nextCatID   int
	nextRestID  int
	nextFoodID  int
	nextComID   int
}

var store *Storage

func InitStorage() {
	store = &Storage{
		users:       make(map[int]*User),
		categories:  make(map[int]*Category),
		restaurants: make(map[int]*Restaurant),
		foods:       make(map[int]*Food),
		comments:    make(map[int]*Comment),
		nextUserID:  1,
		nextCatID:   1,
		nextRestID:  1,
		nextFoodID:  1,
		nextComID:   1,
	}
}

func SeedData() {
	// Create users
	store.users[1] = &User{ID: 1, Name: "John Doe", Email: "john@example.com", CreatedAt: time.Now()}
	store.users[2] = &User{ID: 2, Name: "Jane Smith", Email: "jane@example.com", CreatedAt: time.Now()}
	store.nextUserID = 3

	// Create categories
	store.categories[1] = &Category{ID: 1, Name: "Vietnamese", Description: "Traditional Vietnamese cuisine"}
	store.categories[2] = &Category{ID: 2, Name: "Fast Food", Description: "Quick service meals"}
	store.nextCatID = 3

	// Create restaurants
	store.restaurants[1] = &Restaurant{ID: 1, Name: "Pho 24", Address: "123 Main St", Phone: "0901234567", OwnerID: 1, CreatedAt: time.Now()}
	store.restaurants[2] = &Restaurant{ID: 2, Name: "Banh Mi Saigon", Address: "456 Second St", Phone: "0909876543", OwnerID: 2, CreatedAt: time.Now()}
	store.nextRestID = 3

	// Create foods
	store.foods[1] = &Food{ID: 1, Name: "Pho Bo", Description: "Beef noodle soup", Price: 50000, CategoryID: 1, RestaurantID: 1}
	store.foods[2] = &Food{ID: 2, Name: "Banh Mi Thit", Description: "Pork sandwich", Price: 25000, CategoryID: 2, RestaurantID: 2}
	store.nextFoodID = 3

	// Create comments
	store.comments[1] = &Comment{ID: 1, Content: "Delicious pho!", Rating: 5, UserID: 2, FoodID: 1, CreatedAt: time.Now()}
	store.comments[2] = &Comment{ID: 2, Content: "Great banh mi", Rating: 4, UserID: 1, FoodID: 2, CreatedAt: time.Now()}
	store.nextComID = 3
}
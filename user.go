package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	user.ID = store.nextUserID
	store.nextUserID++
	user.CreatedAt = time.Now()
	store.users[user.ID] = &user
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	users := make([]*User, 0, len(store.users))
	for _, u := range store.users {
		users = append(users, u)
	}
	store.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	store.mu.RLock()
	user, exists := store.users[id]
	store.mu.RUnlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	store.mu.Lock()
	user, exists := store.users[id]
	if !exists {
		store.mu.Unlock()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		store.mu.Unlock()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Giữ nguyên ID và CreatedAt
	updatedUser.ID = user.ID
	updatedUser.CreatedAt = user.CreatedAt
	store.users[id] = &updatedUser
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	store.mu.Lock()
	_, exists := store.users[id]
	if !exists {
		store.mu.Unlock()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	delete(store.users, id)
	store.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func GetUserByID(id int) (*User, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	user, exists := store.users[id]
	return user, exists
}
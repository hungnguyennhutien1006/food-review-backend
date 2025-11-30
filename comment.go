package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Rating    int       `json:"rating"`
	UserID    int       `json:"user_id"`
	FoodID    int       `json:"food_id"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if comment.Rating < 1 || comment.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	_, userExists := GetUserByID(comment.UserID)
	if !userExists {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	_, foodExists := GetFoodByID(comment.FoodID)
	if !foodExists {
		http.Error(w, "Food not found", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	comment.ID = store.nextComID
	store.nextComID++
	comment.CreatedAt = time.Now()
	store.comments[comment.ID] = &comment
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	store.mu.RLock()
	comment, exists := store.comments[id]
	store.mu.RUnlock()

	if !exists {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	store.mu.Lock()
	comment, exists := store.comments[id]
	if !exists {
		store.mu.Unlock()
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	var updatedComment Comment
	if err := json.NewDecoder(r.Body).Decode(&updatedComment); err != nil {
		store.mu.Unlock()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedComment.Rating < 1 || updatedComment.Rating > 5 {
		store.mu.Unlock()
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	// Validate foreign keys
	_, userExists := GetUserByID(updatedComment.UserID)
	_, foodExists := GetFoodByID(updatedComment.FoodID)
	if !userExists || !foodExists {
		store.mu.Unlock()
		http.Error(w, "Invalid user or food", http.StatusBadRequest)
		return
	}

	updatedComment.ID = comment.ID
	updatedComment.CreatedAt = comment.CreatedAt
	store.comments[id] = &updatedComment
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedComment)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	store.mu.Lock()
	_, exists := store.comments[id]
	if !exists {
		store.mu.Unlock()
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	delete(store.comments, id)
	store.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func GetCommentsByFood(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	foodID, _ := strconv.Atoi(vars["food_id"])

	store.mu.RLock()
	comments := make([]*Comment, 0)
	for _, c := range store.comments {
		if c.FoodID == foodID {
			comments = append(comments, c)
		}
	}
	store.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func GetCommentsByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["user_id"])

	store.mu.RLock()
	comments := make([]*Comment, 0)
	for _, c := range store.comments {
		if c.UserID == userID {
			comments = append(comments, c)
		}
	}
	store.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
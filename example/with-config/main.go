package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/k-tsurumaki/fuselage"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var users = map[int]*User{
	1: {ID: 1, Name: "Alice"},
	2: {ID: 2, Name: "Bob"},
}
var nextID = 3

func main() {
	// Load configuration
	config, err := fuselage.LoadConfig("config.yaml")
	if err != nil {
		log.Printf("Failed to load config: %v, using defaults", err)
		config = &fuselage.Config{}
		config.Server.Host = "localhost"
		config.Server.Port = 8080
		config.Middleware = []string{"logger", "recover", "timeout"}
	}
	
	router := fuselage.New()
	
	// Define routes
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUser)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	
	server := fuselage.NewServerFromConfig(config, router)
	log.Printf("Server starting on %s", config.Address())
	log.Fatal(server.ListenAndServe())
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	idStr := fuselage.GetParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	user, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	user.ID = nextID
	nextID++
	users[user.ID] = &user
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	idStr := fuselage.GetParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	if _, exists := users[id]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	user.ID = id
	users[id] = &user
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := fuselage.GetParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	if _, exists := users[id]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
}
package main

import (
	"log"
	"net/http"

	"github.com/k-tsurumaki/fuselage"
)

type User struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=2"`
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
		config.Server.Port = 8083
		config.Middleware = []string{"requestid", "logger", "recover", "timeout"}
	}

	router := fuselage.New()

	// Define routes
	_ = router.GET("/users", getUsers)
	_ = router.GET("/users/:id", getUser)
	_ = router.POST("/users", createUser)
	_ = router.PUT("/users/:id", updateUser)
	_ = router.DELETE("/users/:id", deleteUser)

	server := fuselage.NewServerFromConfig(config, router)
	log.Printf("Server starting on %s", config.Address())
	log.Fatal(server.ListenAndServe())
}

func getUsers(c *fuselage.Context) error {
	return c.JSON(http.StatusOK, users)
}

func getUser(c *fuselage.Context) error {
	id, err := c.ParamInt("id")
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid user ID")
	}

	user, exists := users[id]
	if !exists {
		return c.String(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, user)
}

func createUser(c *fuselage.Context) error {
	var user User
	if err := fuselage.Bind(c, &user); err != nil {
		return err
	}

	user.ID = nextID
	nextID++
	users[user.ID] = &user

	return c.JSON(http.StatusCreated, user)
}

func updateUser(c *fuselage.Context) error {
	id, err := c.ParamInt("id")
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid user ID")
	}

	if _, exists := users[id]; !exists {
		return c.String(http.StatusNotFound, "User not found")
	}

	var user User
	if err := fuselage.Bind(c, &user); err != nil {
		return err
	}

	user.ID = id
	users[id] = &user

	return c.JSON(http.StatusOK, user)
}

func deleteUser(c *fuselage.Context) error {
	id, err := c.ParamInt("id")
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid user ID")
	}

	if _, exists := users[id]; !exists {
		return c.String(http.StatusNotFound, "User not found")
	}

	delete(users, id)
	c.Status(http.StatusNoContent)
	return nil
}
package main

import (
	"log"
	"net/http"

	"github.com/k-tsurumaki/fuselage"
	"github.com/k-tsurumaki/fuselage/middleware"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required,min=2"`
}

var users = map[int]*User{
	1: {ID: 1, Name: "Alice"},
	2: {ID: 2, Name: "Bob"},
}
var nextID = 3

func main() {
	router := fuselage.New()

	// Apply middleware manually
	router.Use(middleware.RequestID())
	router.Use(middleware.Recover())
	router.Use(middleware.Timeout())

	// Define routes
	_ = router.GET("/users", getUsers)
	_ = router.GET("/users/:id", getUser)
	_ = router.POST("/users", createUser, middleware.Logger()) // Logger only for POST /users
	_ = router.PUT("/users/:id", updateUser)
	_ = router.DELETE("/users/:id", deleteUser, middleware.CORS()) // CORS only for DELETE /users/:id

	server := fuselage.NewServer(":8082", router)
	log.Println("Server starting on :8082")
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
	c.SetStatus(http.StatusNoContent)
	return nil
}

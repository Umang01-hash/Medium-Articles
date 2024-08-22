package main

import (
	"github.com/gofr-sample-api/handlers"

	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	h := handlers.New()

	app.GET("/users", h.GetAllUsers)
	app.GET("/users/{id}", h.GetUserByID)
	app.POST("/users", h.AddUser)
	app.PUT("/users/{id}", h.UpdateUser)
	app.DELETE("/users/{id}", h.DeleteUser)

	app.Run()
}

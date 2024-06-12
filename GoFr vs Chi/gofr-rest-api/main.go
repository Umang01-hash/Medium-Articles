package main

import (
	"github.com/gofr-rest-api/handlers"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	h := handlers.New()

	app.GET("/books/{id}", h.GetBookByID)
	app.GET("/books", h.GetAllBooks)
	app.POST("/books", h.AddBook)
	app.PUT("/books/{id}", h.UpdateBook)
	app.DELETE("/books/{id}", h.DeleteBook)

	app.Run()
}

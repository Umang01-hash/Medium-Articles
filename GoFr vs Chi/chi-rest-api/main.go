package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/chi-rest-api/handlers"
)

func initDB() *gorm.DB {
	//Build connection string
	dbUri := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable ",
		"localhost", "postgres", "root123", "2006", "customers")

	db, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	db := initDB()
	h := handlers.New(db)

	r := chi.NewRouter()

	r.Get("/books", h.GetAllBooks)
	r.Post("/books", h.AddBook)
	r.Put("/books/{id}", h.UpdateBook)
	r.Delete("/books/{id}", h.DeleteBook)

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatalf("Error starting the server: %v", err)
		return
	}
}

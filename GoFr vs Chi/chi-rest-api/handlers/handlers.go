package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/chi-rest-api/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// New creates a new handler instance
type handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) handler {
	return handler{db}
}

// GetBookByID retrieves a book by its ID from the database
func (h handler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var book models.Book

	result := h.DB.First(&book, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, book)
}

// GetAllBooks retrieves all books from the database
func (h handler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	var books []models.Book

	result := h.DB.Find(&books)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, books)
}

// AddBook adds a new book to the database
func (h handler) AddBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := h.DB.Create(&book)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, book)
}

// UpdateBook updates an existing book in the database
func (h handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var book models.Book

	result := h.DB.First(&book, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book.ID, _ = strconv.Atoi(id) // Ensure ID remains unchanged
	result = h.DB.Save(&book)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, book)
}

// DeleteBook deletes a book from the database
func (h handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var book models.Book

	result := h.DB.First(&book, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	result = h.DB.Delete(&book)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// respondJSON writes the response as JSON with the given status code
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

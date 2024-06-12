package handlers

import (
	"database/sql"
	"errors"
	"strconv"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	"github.com/gofr-rest-api/models"
)

type handler struct{}

func New() handler {
	return handler{}
}

// GetBookByID retrieves a book by its ID from the database
func (h handler) GetBookByID(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	var book models.Book

	row := ctx.SQL.QueryRowContext(ctx, "SELECT id, title, author FROM books WHERE id = $1", id)
	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, http.ErrorEntityNotFound{Name: "book", Value: id}
		}

		return nil, err
	}

	return book, nil
}

// GetAllBooks retrieves all books from the database
func (h handler) GetAllBooks(ctx *gofr.Context) (interface{}, error) {
	var books []models.Book

	rows, err := ctx.SQL.QueryContext(ctx, "SELECT id, title, author FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and append each book to the books slice
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// AddBook adds a new book to the database
func (h handler) AddBook(ctx *gofr.Context) (interface{}, error) {
	var (
		book models.Book
		id   int
	)

	err := ctx.Bind(&book)
	if err != nil {
		return nil, err
	}

	// Insert the new book into the database and return the generated ID
	row := ctx.SQL.QueryRowContext(ctx, "INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id",
		book.Title, book.Author)
	err = row.Scan(&id)
	if err != nil {
		return nil, err
	}

	book.ID = id

	return book, nil
}

// UpdateBook updates an existing book in the database
func (h handler) UpdateBook(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	var book models.Book

	err := ctx.Bind(&book)
	if err != nil {
		return nil, err
	}

	result, err := ctx.SQL.ExecContext(ctx, "UPDATE books SET title = $1, author = $2 WHERE id = $3",
		book.Title, book.Author, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		// Return a 404 error if no rows were affected (i.e., book not found)
		return nil, http.ErrorEntityNotFound{Name: "book", Value: id}
	}

	book.ID, _ = strconv.Atoi(id)
	return book, nil
}

// DeleteBook deletes a book from the database
func (h handler) DeleteBook(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	result, err := ctx.SQL.ExecContext(ctx, "DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, http.ErrorEntityNotFound{Name: "book", Value: id}
	}

	return nil, nil
}

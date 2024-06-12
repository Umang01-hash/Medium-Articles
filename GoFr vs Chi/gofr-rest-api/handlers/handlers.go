package handlers

import (
	"strconv"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	"github.com/gofr-rest-api/models"
)

type handler struct{}

func New() handler {
	return handler{}
}

func (h handler) GetAllBooks(ctx *gofr.Context) (interface{}, error) {
	var books []models.Book

	rows, err := ctx.SQL.QueryContext(ctx, "SELECT id, title, author FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (h handler) AddBook(ctx *gofr.Context) (interface{}, error) {
	var (
		book models.Book
		id   int
	)

	err := ctx.Bind(&book)
	if err != nil {
		return nil, err
	}

	row := ctx.SQL.QueryRowContext(ctx, "INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id",
		book.Title, book.Author)
	err = row.Scan(&id)
	if err != nil {
		return nil, err
	}

	book.ID = id

	return book, nil
}

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
		return nil, http.ErrorEntityNotFound{Name: "book", Value: id}
	}

	book.ID, _ = strconv.Atoi(id)
	return book, nil
}

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

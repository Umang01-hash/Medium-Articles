package handlers

import (
	"database/sql"
	"errors"
	"strconv"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	"github.com/gofr-sample-api/models"
)

type handler struct{}

func New() handler {
	return handler{}
}

// GetUserByID retrieves a user by its ID from the database
func (h handler) GetUserByID(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	var user models.Users

	row := ctx.SQL.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = ?", id)
	err := row.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, http.ErrorEntityNotFound{Name: "user", Value: id}
		}

		return nil, err
	}

	return user, nil
}

// GetAllUsers retrieves all users from the database
func (h handler) GetAllUsers(ctx *gofr.Context) (interface{}, error) {
	var users []models.Users

	rows, err := ctx.SQL.QueryContext(ctx, "SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and append each user to the users slice
	for rows.Next() {
		var user models.Users
		if err := rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// AddUser adds a new user to the database
func (h handler) AddUser(ctx *gofr.Context) (interface{}, error) {
	var user models.Users

	err := ctx.Bind(&user)
	if err != nil {
		return nil, err
	}

	// Insert the new user into the database
	result, err := ctx.SQL.ExecContext(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.Id = int(id)

	return user, nil
}

// UpdateUser updates an existing user in the database
func (h handler) UpdateUser(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	var user models.Users

	err := ctx.Bind(&user)
	if err != nil {
		return nil, err
	}

	result, err := ctx.SQL.ExecContext(ctx, "UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		// Return a 404 error if no rows were affected (i.e., user not found)
		return nil, http.ErrorEntityNotFound{Name: "user", Value: id}
	}

	user.Id, _ = strconv.Atoi(id)
	return user, nil
}

// DeleteUser deletes a user from the database
func (h handler) DeleteUser(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	result, err := ctx.SQL.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, http.ErrorEntityNotFound{Name: "user", Value: id}
	}

	return nil, nil
}

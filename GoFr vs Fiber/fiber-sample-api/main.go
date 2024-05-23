package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gofiber/fiber/v3"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var db *sql.DB

const (
	createQuery = "INSERT INTO products (name, price) VALUES (?, ?)"
	selectQuery = "SELECT id, name, price FROM products WHERE id = ?"
)

func initDB() {
	var err error

	connectionString := "root:password@tcp(localhost:2001)/test?charset=utf8&parseTime=True&loc=Local"
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	initDB()

	defer db.Close()

	app.Post("/product", func(c fiber.Ctx) error {

		var product Product

		err := c.Bind().Body(&product)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse JSON",
			})
		}

		result, err := db.Exec(createQuery, product.Name, product.Price)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("cannot create product : %v", err),
			})
		}

		id, err := result.LastInsertId()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "cannot retrieve product ID",
			})
		}

		product.ID = int(id)

		return c.Status(fiber.StatusCreated).JSON(product)
	})

	app.Get("/product/:id", func(c fiber.Ctx) error {
		var product Product

		id := c.Params("id")

		err := db.QueryRow(selectQuery, id).Scan(&product.ID, &product.Name, &product.Price)
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "product not found",
			})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "cannot retrieve product",
			})
		}

		return c.JSON(product)
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}

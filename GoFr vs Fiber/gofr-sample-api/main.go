package main

import (
	"gofr.dev/pkg/gofr"
)

const (
	createQuery = "INSERT INTO products (name, price) VALUES (?, ?)"
	selectQuery = "SELECT id, name, price FROM products WHERE id = ?"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func main() {
	app := gofr.New()

	app.POST("/product", func(c *gofr.Context) (interface{}, error) {
		var product Product

		err := c.Bind(&product)
		if err != nil {
			return nil, err
		}

		result, err := c.SQL.Exec(createQuery, product.Name, product.Price)
		if err != nil {
			return nil, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		product.ID = int(id)

		return product, nil
	})

	app.GET("/product/{id}", func(c *gofr.Context) (interface{}, error) {
		var product Product

		id := c.PathParam("id")

		err := c.SQL.QueryRow(selectQuery, id).Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}

		return product, nil
	})

	app.Run()
}

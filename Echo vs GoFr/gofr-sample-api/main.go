package main

import (
	"gofr.dev/pkg/gofr"
)

func main() {
	// Create a new instance of GoFr application
	app := gofr.New()

	// Define a route for handling POST requests to "/redis"
	app.POST("/redis", PostHandler)

	// Define a route for handling GET requests to "/redis/{key}"
	app.GET("/redis/{key}", GetHandler)

	// Start the GoFr application
	app.Run()
}

// PostHandler handles the POST requests to "/redis"
func PostHandler(c *gofr.Context) (interface{}, error) {
	// Parse request body into a map
	input := make(map[string]string)
	if err := c.Request.Bind(&input); err != nil {
		c.Errorf("error while binding request body : %v", err)
		return nil, err
	}

	// Store key-value pairs in Redis
	for key, value := range input {
		err := c.Redis.Set(c, key, value, 0).Err()
		if err != nil {
			c.Error(err)
			return nil, err
		}
	}

	// Return success response
	return "Successful", nil
}

// GetHandler handles the GET requests to "/redis/{key}"
func GetHandler(c *gofr.Context) (interface{}, error) {
	// Extract key from the URL parameter
	key := c.PathParam("key")

	// Retrieve value from Redis based on the key
	value, err := c.Redis.Get(c, key).Result()
	if err != nil {
		c.Error(err)
		return nil, err
	}

	// Prepare response JSON
	resp := make(map[string]string)
	resp[key] = value

	// Return JSON response
	return resp, nil
}

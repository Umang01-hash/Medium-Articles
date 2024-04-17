package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
)

var redisClient *redis.Client // Global variable to hold the Redis client instance

func main() {
	// Initialize the Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:2002",
	})

	// Initialize Echo framework
	e := echo.New()

	e.Use(echoprometheus.NewMiddleware("echo/sample-api")) // adds middleware to gather metrics

	go func() {
		metrics := echo.New()                                // this Echo will run on separate port 8081
		metrics.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
		if err := metrics.Start(":2121"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	// Define routes
	e.POST("/redis", PostUser)    // Route to handle POST requests to /redis
	e.GET("/redis/:key", GetUser) // Route to handle GET requests to /redis/:key

	c := jaegertracing.New(e, nil)
	defer c.Close()

	// Start the server
	e.Logger.Fatal(e.Start(":8000")) // Listen and serve on port 8000
}

// PostUser handles the POST requests to /redis
func PostUser(c echo.Context) error {
	// Parse request body into a map
	input := make(map[string]string)
	err := c.Bind(&input)
	if err != nil {
		c.Logger().Errorf("error while binding request body : %v", err)
		return err
	}

	// Store key-value pairs in Redis
	for key, value := range input {
		err := redisClient.Set(context.Background(), key, value, 0).Err()
		if err != nil {
			c.Logger().Error(err)
			return err
		}
	}

	// Return success response
	return c.String(http.StatusCreated, "Successful")
}

// GetUser handles the GET requests to /redis/:key
func GetUser(c echo.Context) error {
	// Extract key from the URL parameter
	key := c.Param("key")

	// Retrieve value from Redis based on the key
	value, err := redisClient.Get(context.Background(), key).Result()
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	// Prepare response JSON
	resp := make(map[string]string)
	resp[key] = value

	// Return JSON response
	return c.JSON(http.StatusOK, resp)
}

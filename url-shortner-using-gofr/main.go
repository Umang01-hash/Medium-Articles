package main

import (
	"gofr.dev/pkg/gofr"
	"gofr/url-shortner/handlers"
)

func main() {
	app := gofr.New()

	app.POST("/shorten", handlers.ShortenURLHandler)

	// Define route for redirecting shortened URL
	app.GET("/{url}", handlers.RedirectHandler)

	app.Run()
}

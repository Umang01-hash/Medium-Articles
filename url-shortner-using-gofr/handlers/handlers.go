package handlers

import (
	"errors"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"gofr.dev/pkg/gofr"

	"gofr/url-shortner/models"
)

const (
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	expiryHours = 24
)

// base62Encode converts a number to a base62-encoded string
func base62Encode(number uint64) string {
	length := len(alphabet)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(10)
	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(alphabet[(number % uint64(length))])
	}

	return encodedBuilder.String()
}

// ShortenURLHandler handles the shortening of URLs
func ShortenURLHandler(ctx *gofr.Context) (interface{}, error) {
	var (
		req models.Response
		id  string
	)

	err := ctx.Bind(&req)
	if err != nil {
		ctx.Logger.Errorf("error in binding the request. err : %v", err)
		return nil, err
	}

	// Parse URL from request
	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil {
		ctx.Logger.Errorf("malformed or invalid url in request.")
		return nil, err
	}

	// Generate custom short code or use provided one
	if req.CustomShort == "" {
		id = base62Encode(rand.Uint64())
	} else {
		id = req.CustomShort
	}

	// Set default expiry if not provided
	if req.Expiry == 0 {
		req.Expiry = expiryHours
	}

	// Store URL and short code in Redis with expiry
	err = ctx.Redis.Set(ctx, id, parsedURL, req.Expiry*3600*time.Second).Err()
	if err != nil {
		ctx.Logger.Errorf("error while inserting data into redis. err : %v", err)
		return nil, err
	}

	// Populate response object
	resp := models.Response{
		URL:         req.URL,
		CustomShort: id,
		Expiry:      req.Expiry,
	}

	return resp, nil
}

// RedirectHandler handles the redirection of short URLs
func RedirectHandler(ctx *gofr.Context) (interface{}, error) {
	inputURL := ctx.PathParam("url")

	// Retrieve URL from Redis based on short code
	value, err := ctx.Redis.Get(ctx, inputURL).Result()
	switch {
	case errors.Is(err, redis.Nil):
		ctx.Logger.Errorf("short URL not found in database")
		return nil, errors.New("short URL not found in database")
	case err != nil:
		ctx.Logger.Errorf("unable to fetch the URL from Redis: %v", err)
		return nil, err
	}

	return &models.RedirectResponse{OriginalURL: value}, nil
}

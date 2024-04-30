package models

import "time"

type Response struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type RedirectResponse struct {
	OriginalURL string `json:"original_url"`
}

package models

import (
	"encoding/json"
	"time"
)

type Response struct {
	URL         string `json:"url"`
	CustomShort string `json:"short"`
	Expiry      Expiry `json:"expiry"`
}

type RedirectResponse struct {
	OriginalURL string `json:"original_url"`
}

type Expiry time.Duration

// UnmarshalJSON method for Expiry to handle string to time.Duration conversion
func (e *Expiry) UnmarshalJSON(data []byte) error {
	var expiryStr string
	if err := json.Unmarshal(data, &expiryStr); err != nil {
		return err
	}

	expiryDuration, err := time.ParseDuration(expiryStr + "h") // Assuming expiry is in hours
	if err != nil {
		return err
	}

	*e = Expiry(expiryDuration)
	return nil
}

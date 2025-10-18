// Defines structs for JSON responses
package models

// ConversionResponse represents the response structure for currency conversion
type ConversionResponse struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
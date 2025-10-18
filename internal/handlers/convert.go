// Defines the /convert endpoint logic
package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "crypto-project/internal/services" // will hold Coinbase API logic
    "crypto-project/internal/models"
)

// ConvertHandler handles requests to /convert?from=BTC&to=USD
func ConvertHandler(w http.ResponseWriter, r *http.Request) {
    // Read query params
    from := r.URL.Query().Get("from")
    to := r.URL.Query().Get("to")

    if from == "" || to == "" {
        http.Error(w, "Missing 'from' or 'to' query parameter", http.StatusBadRequest)
        return
    }

    // Call service function to get conversion rate
    rate, err := services.GetConversionRate(from, to)
    if err != nil {
        log.Printf("Error fetching conversion rate from %s to %s: %v", from, to, err)
        errorResponse := models.ErrorResponse{
            Error:   "conversion_failed",
            Message: "Unable to fetch conversion rate",
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorResponse)
        return
    }

    // Build response JSON
    response := models.ConversionResponse{
        From: from,
        To:   to,
        Rate: rate,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

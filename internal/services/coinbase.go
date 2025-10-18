//Handles Coinbase API calls
package services

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type CoinbaseRates struct {
    Data struct {
        Currency string            `json:"currency"`
        Rates    map[string]string `json:"rates"`
    } `json:"data"`
}

func GetConversionRate(from string, to string) (float64, error) {
    url := fmt.Sprintf("https://api.coinbase.com/v2/exchange-rates?currency=%s", from)
    resp, err := http.Get(url)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    var rates CoinbaseRates
    if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
        return 0, err
    }

    rateStr, ok := rates.Data.Rates[to]
    if !ok {
        return 0, fmt.Errorf("rate for %s not found", to)
    }

    var rate float64
    fmt.Sscanf(rateStr, "%f", &rate)
    return rate, nil
}

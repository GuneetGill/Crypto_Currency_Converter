package main

import (
    "log"
    "net/http"
    "github.com/GuneetGill/Crypto_Currency_Converter/internal/handlers"
    "github.com/GuneetGill/Crypto_Currency_Converter/internal/middleware"
)

func main() {
    limiter := middleware.NewRateLimiter(10, 1*60*1e9) // 10 requests per minute

    // Serve static files (CSS, JS if needed) at /static/
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

    // Register the /convert endpoint (with rate limiter)
    http.Handle("/convert", limiter.Limit(http.HandlerFunc(handlers.ConvertHandler)))

    // Serve the main HTML page at "/" (also rate limited)
    http.Handle("/", limiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Serving request for path: %s", r.URL.Path)
        http.ServeFile(w, r, "./web/templates/index.html")
    })))

    // Start the HTTP server on port 8080
    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

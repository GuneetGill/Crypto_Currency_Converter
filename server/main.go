package main

import (
    "log"
    "net/http"
    "crypto-project/internal/handlers"
)

func main() {
    limiter := middleware.NewRateLimiter(10, 1*60*1e9) // 10 requests per minute

    // Serve static files (CSS, JS if needed) at /static/
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

    // Register the /convert endpoint
    http.HandleFunc("/convert", handlers.ConvertHandler)

    // Serve the main HTML page at "/" (this should be last as it's a catch-all)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Serving request for path: %s", r.URL.Path)
        http.ServeFile(w, r, "./web/templates/index.html")
    })

    // Start the HTTP server on port 8080
    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

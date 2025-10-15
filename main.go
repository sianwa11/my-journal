package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sianwa11/my-journal/internal/api/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	routes := routes.SetupRoutes()

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        routes,
		ReadTimeout:    15 * time.Second, // Time to read request
		WriteTimeout:   15 * time.Second, // Time to write response
		IdleTimeout:    60 * time.Second, // Time to keep connection alive
		MaxHeaderBytes: 1 << 20,          // 1 MB max header size
	}

	log.Printf("Serving on: http://localhost:%s/api/\n", port)
	log.Fatal(server.ListenAndServe())
}

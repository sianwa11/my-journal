package main

import (
	"log"
	"net/http"

	"github.com/sianwa11/my-journal/internal/api/routes"
)

func main() {

	routes := routes.SetupRoutes()
	port := "8080"

	log.Printf("Serving on: http://localhost:%s/api/\n", port)
	http.ListenAndServe(":"+port, routes)
}

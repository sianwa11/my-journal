package main

import (
	"net/http"

	"github.com/sianwa11/my-journal/internal/api/routes"
)


func main() {

	routes := routes.SetupRoutes()

	http.ListenAndServe(":8080", routes)
}

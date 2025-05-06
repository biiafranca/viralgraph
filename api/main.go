package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/biiafranca/viralgraph/api/routes"
)

func main() {
	r := chi.NewRouter()

	routes.RegisterCovidStatsRoutes(r)
	//routes.RegisterVaccineRoutes(r)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

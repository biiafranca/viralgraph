package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/biiafranca/viralgraph/api/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	routes.RegisterCovidStatsRoutes(r)
	routes.RegisterVaccinationRoutes(r)
	routes.RegisterUsedVaccinesRoutes(r)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

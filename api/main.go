package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/biiafranca/viralgraph/api/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	routes.RegisterCovidStatsRoutes(r)
	routes.RegisterVaccinationRoutes(r)
	routes.RegisterUsedVaccinesRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server listening on 0.0.0.0:%s\n", port)
	fmt.Println("If you're running locally, access: http://localhost:" + port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error to start server: %v\n", err)
		os.Exit(1)
	}
}

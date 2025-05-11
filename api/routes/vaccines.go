// Package routes defines the application's URL routing.
// This file registers all routes related to COVID-19 vaccines,
// and connects each endpoint to its corresponding handler.
//
// Specifically, it defines routes for the /vaccines endpoints.

package routes

import (
	"github.com/biiafranca/viralgraph/api/handlers/vaccines"
	"github.com/go-chi/chi/v5"
)

func RegisterUsedVaccinesRoutes(r chi.Router) {
	r.Get("/vaccines", vaccines.HandleVaccines)
	r.Get("/vaccines/used-in/{country}", vaccines.HandleUsedInCountry)
	r.Get("/vaccines/first-use", vaccines.HandleFirstUse)
	r.Get("/vaccines/{vaccineID}/used-by", vaccines.HandleUsedBy)
}

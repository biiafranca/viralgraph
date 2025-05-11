// Package routes defines the application's URL routing.
// This file registers all routes related to COVID-19 vaccination,
// and connects each endpoint to its corresponding controller.
//
// Specifically, it defines routes for the /vaccination endpoints.

package routes

import (
	"github.com/biiafranca/viralgraph/api/handlers/vaccination"
	"github.com/go-chi/chi/v5"
)

func RegisterVaccinationRoutes(r chi.Router) {
	r.Get("/vaccination/{country}/{date}", vaccination.VaccinationController)
	r.Get("/vaccination/{date}", vaccination.VaccinationController)
}

package routes

import (
	"github.com/biiafranca/viralgraph/api/handlers/vaccination"
	"github.com/go-chi/chi/v5"
)

func RegisterVaccinationRoutes(r chi.Router) {
	r.Get("/vaccination/{country}/{date}", vaccination.VaccinationController)
	r.Get("/vaccination/{date}", vaccination.VaccinationController)
}

package routes

import (
	"github.com/biiafranca/viralgraph/api/handlers/usedvaccines"
	"github.com/go-chi/chi/v5"
)

func RegisterUsedVaccinesRoutes(r chi.Router) {
	r.Get("/used-vaccines/{country}", usedvaccines.HandleUsedVaccines)
}

// Package routes defines the application's URL routing.
// This file registers all routes related to COVID-19 statistics,
// and connects each endpoint to its corresponding controller.
//
// Specifically, it defines routes for the /covid-stats endpoints.

package routes

import (
	"github.com/biiafranca/viralgraph/api/handlers/covidstats"
	"github.com/go-chi/chi/v5"
)

func RegisterCovidStatsRoutes(r chi.Router) {
	// Local stats, by country and date (ex: /covid-stats/BRA/2021-01-01)
	r.Get("/covid-stats/{country}/{date}", covidstats.CovidStatsController)

	// Global stats, by date (ex: /covid-stats/2021-01-01)
	r.Get("/covid-stats/{date}", covidstats.CovidStatsController)
}

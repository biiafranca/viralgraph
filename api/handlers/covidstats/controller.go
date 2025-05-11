// Package covidstats handles COVID-19 case statistics.
// This file routes incoming HTTP requests to the appropriate handler based on the URL and query parameters.
//
// It supports both accumulated and daily data, at country or global level.

package covidstats

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func CovidStatsController(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")
	date := chi.URLParam(r, "date")
	onlyNews := strings.ToLower(r.URL.Query().Get("only-news")) == "true"

	if onlyNews {
		handleNew(w, country, date)
	} else {
		handleAccumulated(w, country, date)
	}
}

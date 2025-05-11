package vaccination

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func VaccinationController(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")
	date := chi.URLParam(r, "date")
	onlyNews := strings.ToLower(r.URL.Query().Get("only-news")) == "true"

	if onlyNews {
		handleNew(w, country, date)
	} else {
		handleAccumulated(w, country, date)
	}
}

package vaccination

import (
	"net/http"
	"strings"
)

func VaccinationController(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/vaccination/")
	parts := strings.Split(path, "/")

	var country, date string
	if len(parts) == 1 {
		date = parts[0]
	} else if len(parts) == 2 {
		country = parts[0]
		date = parts[1]
	} else {
		http.NotFound(w, r)
		return
	}

	onlyNews := strings.ToLower(r.URL.Query().Get("only-news")) == "true"

	if onlyNews {
		handleNew(w, country, date)
	} else {
		handleAccumulated(w, country, date)
	}
}

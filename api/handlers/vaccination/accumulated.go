package vaccination

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/biiafranca/viralgraph/api/neo4j"
	"github.com/biiafranca/viralgraph/api/utils"
)

func handleAccumulated(w http.ResponseWriter, country, date string) {

	// Validate date:
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD.")
		return
	}
	if parsedDate.After(time.Now()) {
		utils.RespondWithError(w, http.StatusNotFound, "No data available for future dates")
		return
	}

	ctx := context.Background()
	session := neo4j.GetSession()
	defer session.Close(ctx)

	var query string
	var params map[string]interface{}

	if country != "" {
		query = `
			MATCH (c:Country {iso3: $country})-[:VACCINATED_ON]->(vs:VaccinationStats)
			WHERE vs.date <= date($date)
			RETURN vs.totalVaccinated AS totalVaccinated
			ORDER BY vs.date DESC
			LIMIT 1
		`
		params = map[string]interface{}{"country": country, "date": date}
	} else {
		query = `
			MATCH (c:Country)-[:VACCINATED_ON]->(vs:VaccinationStats)
			WHERE vs.date <= date($date)
			WITH c, vs ORDER BY vs.date DESC
			WITH c, collect(vs)[0] AS latest
			RETURN sum(latest.totalVaccinated) AS totalVaccinated
		`
		params = map[string]interface{}{"date": date}
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		log.Printf("Neo4j query failed: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to query database")
		return
	}

	if !result.Next(ctx) {
		utils.RespondWithError(w, http.StatusNotFound, "No data found for the given input")
		return
	}

	record := result.Record()
	totalVaccinatedRaw, _ := record.Get("totalVaccinated")
	totalVaccinated := totalVaccinatedRaw.(int64)

	label := "worldwide"
	if country != "" {
		label = country
	}

	response := VaccinationResponse{
		Country:         label,
		Date:            date,
		OnlyNews:        false,
		TotalVaccinated: totalVaccinated,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Package vaccination handles COVID-19 vaccination statistics.
// This file implements logic for calculating *new daily vaccinations*.
//
// The handler defined here is used when the `onlyNews` parameter is true.
// It returns the difference between the requested date and the most recent prior date,
// effectively showing the new vaccinations for the given day.

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

func handleNew(w http.ResponseWriter, country, date string) {

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

	var currentQuery, previousQuery string
	var params map[string]interface{}

	if country != "" {
		currentQuery = `
			MATCH (c:Country {iso3: $country})-[:VACCINATED_ON]->(vs:VaccinationStats)
			WHERE vs.date = date($date)
			RETURN vs.totalVaccinated AS totalVaccinated
		`
		previousQuery = `
			MATCH (c:Country {iso3: $country})-[:VACCINATED_ON]->(vs:VaccinationStats)
			WHERE vs.date < date($date)
			RETURN vs.totalVaccinated AS totalVaccinated
			ORDER BY vs.date DESC
			LIMIT 1
		`
		params = map[string]interface{}{"country": country, "date": date}
	} else {
		currentQuery = `
			MATCH (c:Country)-[:VACCINATED_ON]->(vs:VaccinationStats)
			WHERE vs.date = date($date)
			RETURN sum(vs.totalVaccinated) AS totalVaccinated
		`
		previousQuery = `
			MATCH (c:Country)-[:VACCINATED_ON]->(cur:VaccinationStats)
			WHERE cur.date = date($date)
			WITH c

			MATCH (c)-[:VACCINATED_ON]->(prev:VaccinationStats)
			WHERE prev.date < date($date)
			WITH c, prev ORDER BY prev.date DESC
			WITH c, collect(prev)[0] AS latest
			RETURN sum(latest.totalVaccinated) AS totalVaccinated
		`
		params = map[string]interface{}{"date": date}
	}

	currentRes, err := session.Run(ctx, currentQuery, params)
	if err != nil {
		log.Printf("Neo4j query failed: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to query current date data")
		return
	}
	if !currentRes.Next(ctx) {
		utils.RespondWithError(w, http.StatusNotFound, "No data found for the current date")
		return
	}

	record := currentRes.Record()
	currentVaccinatedRaw, _ := record.Get("totalVaccinated")
	currentVaccinated := currentVaccinatedRaw.(int64)

	previousVaccinated := int64(0)
	prevRes, err := session.Run(ctx, previousQuery, params)
	if err == nil && prevRes.Next(ctx) {
		record := prevRes.Record()
		previousVaccinatedRaw, _ := record.Get("totalVaccinated")
		previousVaccinated = previousVaccinatedRaw.(int64)
	}

	newVaccinated := currentVaccinated - previousVaccinated

	label := "worldwide"
	if country != "" {
		label = country
	}

	response := VaccinationResponse{
		Country:         label,
		Date:            date,
		OnlyNews:        true,
		TotalVaccinated: newVaccinated,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

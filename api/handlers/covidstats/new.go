// Package covidstats handles COVID-19 case statistics.
// This file implements logic for calculating *new daily cases and deaths*.
//
// The handler defined here is used when the `onlyNews` parameter is true.
// It returns the difference between the requested date and the most recent prior date,
// effectively showing the new cases and deaths for the given day.

package covidstats

import (
	"context"
	"encoding/json"
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
			MATCH (c:Country {iso3: $country})-[:HAS_CASE]->(cc:CovidCase)
			WHERE cc.date = date($date)
			RETURN cc.totalCases AS totalCases, cc.totalDeaths AS totalDeaths
		`
		previousQuery = `
			MATCH (c:Country {iso3: $country})-[:HAS_CASE]->(cc:CovidCase)
			WHERE cc.date < date($date)
			RETURN cc.totalCases AS totalCases, cc.totalDeaths AS totalDeaths
			ORDER BY cc.date DESC
			LIMIT 1
		`
		params = map[string]interface{}{"country": country, "date": date}
	} else {
		currentQuery = `
			MATCH (c:Country)-[:HAS_CASE]->(cc:CovidCase)
			WHERE cc.date = date($date)
			RETURN sum(cc.totalCases) AS totalCases, sum(cc.totalDeaths) AS totalDeaths
		`
		previousQuery = `
			MATCH (c:Country)-[:HAS_CASE]->(cc:CovidCase)
			WHERE cc.date < date($date)
			WITH c, cc ORDER BY cc.date DESC
			WITH c, collect(cc)[0] AS latest
			RETURN sum(latest.totalCases) AS totalCases, sum(latest.totalDeaths) AS totalDeaths
		`
		params = map[string]interface{}{"date": date}
	}

	currentRes, err := session.Run(ctx, currentQuery, params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to query current date data")
		return
	}
	if !currentRes.Next(ctx) {
		utils.RespondWithError(w, http.StatusNotFound, "No data found for the current date")
		return
	}

	record := currentRes.Record()
	currentCasesRaw, _ := record.Get("totalCases")
	currentDeathsRaw, _ := record.Get("totalDeaths")
	currentCases := currentCasesRaw.(int64)
	currentDeaths := currentDeathsRaw.(int64)

	previousCases, previousDeaths := int64(0), int64(0)
	prevRes, err := session.Run(ctx, previousQuery, params)
	if err == nil && prevRes.Next(ctx) {
		record := prevRes.Record()
		prevCasesRaw, _ := record.Get("totalCases")
		prevDeathsRaw, _ := record.Get("totalDeaths")
		previousCases = prevCasesRaw.(int64)
		previousDeaths = prevDeathsRaw.(int64)
	}

	newCases := currentCases - previousCases
	newDeaths := currentDeaths - previousDeaths

	label := "worldwide"
	if country != "" {
		label = country
	}

	response := CovidStatsResponse{
		Country:  label,
		Date:     date,
		OnlyNews: true,
		Cases:    newCases,
		Deaths:   newDeaths,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

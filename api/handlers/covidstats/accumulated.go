// Package covidstats handles COVID-19 case statistics.
// This file contains the logic for calculating *accumulated* cases and deaths.
//
// It is used when the `onlyNews` parameter is false or absent.
// The total reflects the last known values *on or before* the requested date.

package covidstats

import (
	"context"
	"encoding/json"
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
			MATCH (c:Country {iso3: $country})-[:HAS_CASE]->(cc:CovidCase)
			WHERE cc.date <= date($date)
			RETURN cc.totalCases AS totalCases, cc.totalDeaths AS totalDeaths
			ORDER BY cc.date DESC
			LIMIT 1
		`
		params = map[string]interface{}{"country": country, "date": date}
	} else {
		query = `
			MATCH (c:Country)-[:HAS_CASE]->(cc:CovidCase)
			WHERE cc.date <= date($date)
			WITH c, cc ORDER BY cc.date DESC
			WITH c, collect(cc)[0] AS latest
			RETURN sum(latest.totalCases) AS totalCases, sum(latest.totalDeaths) AS totalDeaths
		`
		params = map[string]interface{}{"date": date}
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to query database")
		return
	}

	if !result.Next(ctx) {
		utils.RespondWithError(w, http.StatusNotFound, "No data found for the given input")
		return
	}

	record := result.Record()
	totalCases, _ := record.Get("totalCases")
	totalDeaths, _ := record.Get("totalDeaths")
	
	label := "worldwide"
	if country != "" {
		label = country
	}

	response := CovidStatsResponse{
		Country:  label,
		Date:     date,
		OnlyNews: false,
		Cases:    totalCases.(int64),
		Deaths:   totalDeaths.(int64),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

package vaccines

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/biiafranca/viralgraph/api/neo4j"
	"github.com/biiafranca/viralgraph/api/utils"
)

func HandleUsedBy(w http.ResponseWriter, r *http.Request) {
	vaccineIDStr := r.PathValue("vaccineID")
	if vaccineIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Vaccine ID parameter is required")
		return
	}
	vaccineID, err := strconv.Atoi(vaccineIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Vaccine ID must be an integer")
		return
	}

	ctx := context.Background()
	session := neo4j.GetSession()
	defer session.Close(ctx)

	queryName := `
		MATCH (v:Vaccine {id: $id_vaccine})
		RETURN v.name As name
	`
	query := `
		MATCH (c:Country)-[r:USES]->(v:Vaccine {id: $id_vaccine})
		RETURN c.iso3 AS country, r.first_used AS date
		ORDER BY country
	`
	params := map[string]interface{}{"id_vaccine": vaccineID}

	nameRes, errName := session.Run(ctx, queryName, params)
	if errName != nil {
		log.Printf("Neo4j query failed: %v", errName)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to query database")
		return
	}
	if !nameRes.Next(ctx) {
		utils.RespondWithError(w, http.StatusNotFound, "No vaccine found for this ID")
		return
	}

	record := nameRes.Record()
	contextName, _ := record.Get("name")
	contextNameStr := contextName.(string)

	result, err := session.Run(ctx, query, params)
	if err != nil {
		log.Printf("Neo4j query failed: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to query database")
		return
	}

	var countryUsage []UsageEntry
	for result.Next(ctx) {
		record := result.Record()
		country, _ := record.Get("country")
		date, _ := record.Get("date")
		countryUsage = append(countryUsage, UsageEntry{
			Country:  country.(string),
			FirstUse: fmt.Sprint(date),
		})
	}

	if len(countryUsage) == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "No usage data found for this vaccine")
		return
	}

	response := UsageResponse{
		Context: contextNameStr,
		Entries: countryUsage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

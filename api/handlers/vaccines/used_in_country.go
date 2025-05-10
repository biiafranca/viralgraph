package vaccines

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/biiafranca/viralgraph/api/neo4j"
	"github.com/biiafranca/viralgraph/api/utils"
)

func HandleUsedInCountry(w http.ResponseWriter, r *http.Request) {

	country := strings.ToUpper(r.PathValue("country"))
	if country == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Country parameter is required")
		return
	}

	ctx := context.Background()
	session := neo4j.GetSession()
	defer session.Close(ctx)

	query := `
		MATCH (c:Country {iso3: $country})-[r:USES]->(v:Vaccine)
		RETURN v.name AS vaccine, r.first_used AS date
		ORDER BY vaccine
	`

	params := map[string]interface{}{"country": country}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		log.Printf("Neo4j query failed: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to query database")
		return
	}

	var vaccineUsage []UsageEntry
	for result.Next(ctx) {
		record := result.Record()
		vaccine, _ := record.Get("vaccine")
		date, _ := record.Get("date")
		vaccineUsage = append(vaccineUsage, UsageEntry{
			Vaccine:  vaccine.(string),
			FirstUse: fmt.Sprint(date),
		})
	}

	if len(vaccineUsage) == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "No vaccines found for this country")
		return
	}

	response := UsageResponse{
		Context: country,
		Entries: vaccineUsage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

package vaccines

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/biiafranca/viralgraph/api/neo4j"
	"github.com/biiafranca/viralgraph/api/utils"
)

func HandleFirstUse(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	session := neo4j.GetSession()
	defer session.Close(ctx)

	query := `
		MATCH (v:Vaccine)
		WHERE v.first_global_use IS NOT NULL
		RETURN v.name AS vaccine, v.first_global_use AS date
		ORDER BY date
	`

	result, err := session.Run(ctx, query, nil)
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

	response := UsageResponse{
		Context: "worldwide",
		Entries: vaccineUsage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

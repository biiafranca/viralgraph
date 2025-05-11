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

func HandleVaccines(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	session := neo4j.GetSession()
	defer session.Close(ctx)

	query := `
		MATCH (v:Vaccine)
		RETURN v.id AS id, v.name AS vaccine, v.first_global_use AS date
		ORDER BY id
	`

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		log.Printf("Neo4j query failed: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to query database")
		return
	}

	var vaccines []Vaccine
	for result.Next(ctx) {
		record := result.Record()
		id, _ := record.Get("id")
		vaccine, _ := record.Get("vaccine")
		date, _ := record.Get("date")
		vaccines = append(vaccines, Vaccine{
			ID:             int(id.(int64)),
			Name:           vaccine.(string),
			FirstGlobalUse: fmt.Sprint(date),
		})
	}

	response := VaccinesResponse{
		Vaccines: vaccines,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

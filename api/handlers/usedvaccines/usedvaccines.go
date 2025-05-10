package usedvaccines

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/biiafranca/viralgraph/api/neo4j"
	"github.com/biiafranca/viralgraph/api/utils"
)

func HandleUsedVaccines(w http.ResponseWriter, r *http.Request) {

	country := strings.ToUpper(r.PathValue("country"))
	if country == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Country parameter is required")
		return
	}

	ctx := context.Background()
	session := neo4j.GetSession()
	defer session.Close(ctx)

	query := `
		MATCH (c:Country {iso3: $country})-[:USES]->(v:Vaccine)
		RETURN v.name AS vaccine
		ORDER BY vaccine
	`

	params := map[string]interface{}{"country": country}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to query database")
		return
	}

	var vaccines []string
	for result.Next(ctx) {
		record := result.Record()
		vaccine, _ := record.Get("vaccine")
		vaccines = append(vaccines, vaccine.(string))
	}

	if len(vaccines) == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "No vaccines found for this country")
		return
	}

	response := UsedVaccinesResponse{
		Country:      country,
		UsedVaccines: vaccines,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

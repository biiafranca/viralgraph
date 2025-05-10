package usedvaccines

type UsedVaccinesResponse struct {
	Country      string   `json:"country"`
	UsedVaccines []string `json:"used_vaccines"`
}

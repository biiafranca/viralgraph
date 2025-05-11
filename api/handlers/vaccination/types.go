package vaccination

type VaccinationResponse struct {
	Country         string `json:"country"`
	Date            string `json:"date"`
	OnlyNews        bool   `json:"only_news"`
	TotalVaccinated int64  `json:"total_vaccinated"`
}

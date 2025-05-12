// Package vaccination handles COVID-19 vaccination statistics.
// Defines response data structures used across vaccination handlers.

package vaccination

type VaccinationResponse struct {
	Country         string `json:"country"`
	Date            string `json:"date"`
	OnlyNews        bool   `json:"only_news"`
	TotalVaccinated int64  `json:"total_vaccinated"`
}

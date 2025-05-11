// Package vaccines handles COVID-19 vaccines information.
// Defines response data structures used across vaccines handlers.
//
// Usage-related types (UsageEntry, UsageResponse) are used to represent generic
// statistics of vaccine usage, either by country or by vaccine.
//
// Vaccine-related types (Vaccine, VaccinesResponse) focus on providing general
// metadata about each vaccine.

package vaccines

type UsageEntry struct {
	Country  string `json:"country,omitempty"`
	Vaccine  string `json:"vaccine,omitempty"`
	FirstUse string `json:"first_use"`
}

type UsageResponse struct {
	Context string       `json:"context"`
	Entries []UsageEntry `json:"entries"`
}

type Vaccine struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	FirstGlobalUse string `json:"first_global_use"`
}

type VaccinesResponse struct {
	Vaccines []Vaccine `json:"vaccines"`
}

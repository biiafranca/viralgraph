// Package covidstats handles COVID-19 case statistics.
// Defines response data structures used across COVID-19 statistics handlers.

package covidstats

type CovidStatsResponse struct {
	Country  string `json:"country"`
	Date     string `json:"date"`
	OnlyNews bool   `json:"only_news"`
	Cases    int64  `json:"cases"`
	Deaths   int64  `json:"deaths"`
}

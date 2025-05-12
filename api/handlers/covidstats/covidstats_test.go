package covidstats

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestHandleNew_InvalidDate(t *testing.T) {
	rec := httptest.NewRecorder()
	handleNew(rec, "", "invalid-date")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleNew_FutureDate(t *testing.T) {
	future := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	rec := httptest.NewRecorder()
	handleNew(rec, "", future)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d for future date, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHandleNew_PositiveGlobal(t *testing.T) {
	// Date present in the test database
	date := "2021-07-31"
	rec := httptest.NewRecorder()
	handleNew(rec, "", date)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d for valid global stats, got %d", http.StatusOK, rec.Code)
	}
}

func TestHandleNew_PositiveCountry(t *testing.T) {
	// Real ISO3 code and date present in the test database
	country := "BRA"
	date := "2021-07-31"
	rec := httptest.NewRecorder()
	handleNew(rec, country, date)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d for valid country stats, got %d", http.StatusOK, rec.Code)
	}
}

func TestHandleAccumulated_InvalidDate(t *testing.T) {
	rec := httptest.NewRecorder()
	handleAccumulated(rec, "", "bad-date")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for invalid date, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleAccumulated_PositiveGlobal(t *testing.T) {
	// Date present in the test database
	date := "2021-07-31"
	rec := httptest.NewRecorder()
	handleAccumulated(rec, "", date)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d for global accumulated, got %d", http.StatusOK, rec.Code)
	}
}

func TestHandleAccumulated_PositiveCountry(t *testing.T) {
	// Real ISO3 code and date present in the test database
	country := "BRA"
	date := "2021-07-31"
	rec := httptest.NewRecorder()
	handleAccumulated(rec, country, date)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d for country accumulated, got %d", http.StatusOK, rec.Code)
	}
}

func TestCovidStatsController_Routes(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/covid-stats/{country}/{date}", CovidStatsController)
	r.Get("/covid-stats/{date}", CovidStatsController)

	// Test global route
	req1 := httptest.NewRequest(http.MethodGet, "/covid-stats/2021-07-31?onlyNews=false", nil)
	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Errorf("expected status %d for global accumulated, got %d", http.StatusOK, rec1.Code)
	}

	// Test country route with onlyNews=true
	req2 := httptest.NewRequest(http.MethodGet, "/covid-stats/BRA/2021-07-31?onlyNews=true", nil)
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Errorf("expected status %d for country new data, got %d", http.StatusOK, rec2.Code)
	}
}

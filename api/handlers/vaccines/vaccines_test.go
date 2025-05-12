package vaccines

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestHandleFirstUse_Positive(t *testing.T) {
	rec := httptest.NewRecorder()
	HandleFirstUse(rec, httptest.NewRequest(http.MethodGet, "/vaccines/first-use", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestHandleFirstUse_Route(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/vaccines/first-use", HandleFirstUse)

	req := httptest.NewRequest(http.MethodGet, "/vaccines/first-use", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d on /vaccines/first-use route, got %d", http.StatusOK, rec.Code)
	}
}

func TestHandleUsedBy_Route_Positive(t *testing.T) {
	r := chi.NewRouter()
	r.Route("/vaccines", func(r chi.Router) {
		r.Get("/{vaccineID}/used-by", HandleUsedBy)
	})

	req := httptest.NewRequest(http.MethodGet, "/vaccines/1/used-by", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestHandleUsedBy_Route_NonexistentVaccine(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/vaccines/{vaccineID}/used-by", HandleUsedBy)

	req := httptest.NewRequest(http.MethodGet, "/vaccines/99999/used-by", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 for nonexistent vaccine, got %d", rec.Code)
	}
}

func TestHandleUsedInCountry_Route_Positive(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/vaccines/used-in/{country}", HandleUsedInCountry)

	req := httptest.NewRequest(http.MethodGet, "/vaccines/used-in/BRA", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestHandleUsedInCountry_Route_NonexistentCountry(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/vaccines/used-in/{country}", HandleUsedInCountry)

	req := httptest.NewRequest(http.MethodGet, "/vaccines/used-in/XYZ", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 for nonexistent country, got %d", rec.Code)
	}
}

func TestHandleVaccines_Positive(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/vaccines", nil)
	rec := httptest.NewRecorder()
	HandleVaccines(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestHandleVaccines_Route(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/vaccines", HandleVaccines)

	req := httptest.NewRequest(http.MethodGet, "/vaccines", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

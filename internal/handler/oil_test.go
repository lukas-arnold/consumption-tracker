package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func TestOilView(t *testing.T) {
	h := testHandler(t)
	req := httptest.NewRequest("GET", "/oil", nil)
	rec := httptest.NewRecorder()

	h.HandleOilView(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestAddOilGet(t *testing.T) {
	h := testHandler(t)
	req := httptest.NewRequest("GET", "/oil/add", nil)
	rec := httptest.NewRecorder()

	h.HandleAddOilGet(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestAddOilPost(t *testing.T) {
	h := testHandler(t)
	form := url.Values{"date": {"2024-01-01"}, "volume": {"1000"}, "costs": {"900"}, "retailer": {"Shell"}}

	req := httptest.NewRequest("POST", "/oil/add", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.HandleAddOilPost(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	oils, _ := h.store.GetOil()
	if len(oils) != 1 || oils[0].Volume != 1000 {
		t.Fatal("oil entry not created correctly")
	}
}

func TestEditOil(t *testing.T) {
	h := testHandler(t)
	h.store.AddOil(models.OilInput{Date: "2024"})
	oils, _ := h.store.GetOil()

	req := httptest.NewRequest("GET", "/oil/edit", nil)
	req.SetPathValue("id", strconv.FormatInt(oils[0].Id, 10))
	rec := httptest.NewRecorder()

	h.HandleEditOil(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestEditOilInvalidID(t *testing.T) {
	h := testHandler(t)
	req := httptest.NewRequest("GET", "/oil/edit", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleEditOil(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestSaveOil(t *testing.T) {
	h := testHandler(t)
	h.store.AddOil(models.OilInput{Date: "2024", Volume: 100})
	oils, _ := h.store.GetOil()

	form := url.Values{"date": {"2025"}, "volume": {"200"}, "costs": {"500"}}
	req := httptest.NewRequest("POST", "/oil/save", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", strconv.FormatInt(oils[0].Id, 10))
	rec := httptest.NewRecorder()

	h.HandleSaveOil(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	updated, _ := h.store.GetOil()
	if updated[0].Volume != 200 {
		t.Fatal("oil not updated")
	}
}

func TestSaveOilInvalidID(t *testing.T) {
	h := testHandler(t)
	form := url.Values{"date": {"2025-01-01"}, "volume": {"1000"}, "costs": {"900"}}

	req := httptest.NewRequest("POST", "/oil/save", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleSaveOil(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestDeleteOil(t *testing.T) {
	h := testHandler(t)
	h.store.AddOil(models.OilInput{Date: "2024"})
	oils, _ := h.store.GetOil()

	req := httptest.NewRequest("GET", "/oil/delete", nil)
	req.SetPathValue("id", strconv.FormatInt(oils[0].Id, 10))
	rec := httptest.NewRecorder()

	h.HandleDeleteOil(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	result, _ := h.store.GetOil()
	if len(result) != 0 {
		t.Fatal("oil not deleted")
	}
}

func TestDeleteOilInvalidID(t *testing.T) {
	h := testHandler(t)
	req := httptest.NewRequest("GET", "/oil/delete", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleDeleteOil(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestAddOilFillLevelGet(t *testing.T) {
	h := testHandler(t)
	req := httptest.NewRequest("GET", "/oil/fill/add", nil)
	rec := httptest.NewRecorder()

	h.HandleAddOilFillLevelGet(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestAddOilFillLevelPost(t *testing.T) {
	h := testHandler(t)
	form := url.Values{"date": {"2024"}, "level": {"80"}}

	req := httptest.NewRequest("POST", "/oil/fill/add", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.HandleAddOilFillLevelPost(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	levels, _ := h.store.GetOilFillLevels()
	if len(levels) != 1 {
		t.Fatal("fill level not created")
	}
}

func TestEditOilFillLevel(t *testing.T) {
	h := testHandler(t)
	h.store.AddOilFillLevel(models.OilFillLevelInput{Date: "2024-01-01", Level: 80})
	levels, _ := h.store.GetOilFillLevels()

	req := httptest.NewRequest("GET", "/oil/fill/edit", nil)
	req.SetPathValue("id", strconv.FormatInt(levels[0].Id, 10))
	rec := httptest.NewRecorder()

	h.HandleEditOilFillLevel(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestEditOilFillLevelInvalidID(t *testing.T) {
	h := testHandler(t)
	req := httptest.NewRequest("GET", "/oil/fill/edit", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleEditOilFillLevel(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestSaveOilFillLevel(t *testing.T) {
	h := testHandler(t)
	h.store.AddOilFillLevel(models.OilFillLevelInput{Date: "2024", Level: 50})
	levels, _ := h.store.GetOilFillLevels()

	form := url.Values{"date": {"2025"}, "level": {"100"}}
	req := httptest.NewRequest("POST", "/oil/fill/save", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", strconv.FormatInt(levels[0].Id, 10))
	rec := httptest.NewRecorder()

	h.HandleSaveOilFillLevel(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	updated, _ := h.store.GetOilFillLevels()
	if updated[0].Level != 100 || updated[0].Date != "2025" {
		t.Fatal("fill level not updated correctly")
	}
}

func TestSaveOilFillLevelInvalidID(t *testing.T) {
	h := testHandler(t)
	form := url.Values{"date": {"2025-01-01"}, "level": {"80"}}

	req := httptest.NewRequest("POST", "/oil/fill/save", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleSaveOilFillLevel(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestDeleteOilFillLevel(t *testing.T) {
	h := testHandler(t)
	h.store.AddOilFillLevel(models.OilFillLevelInput{Date: "2024", Level: 70})
	levels, _ := h.store.GetOilFillLevels()

	req := httptest.NewRequest("GET", "/oil/fill/delete", nil)
	req.SetPathValue("id", strconv.FormatInt(levels[0].Id, 10))
	rec := httptest.NewRecorder()

	h.HandleDeleteOilFillLevel(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	result, _ := h.store.GetOilFillLevels()
	if len(result) != 0 {
		t.Fatal("fill level not deleted")
	}
}

func TestDeleteOilFillLevelInvalidID(t *testing.T) {
	h := testHandler(t)
	req := httptest.NewRequest("GET", "/oil/fill/delete", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleDeleteOilFillLevel(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

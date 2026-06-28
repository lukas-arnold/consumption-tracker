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

func TestHandleWaterView(t *testing.T) {
	h := testHandler(t)

	req := httptest.NewRequest("GET", "/water", nil)
	rec := httptest.NewRecorder()

	h.HandleWaterView(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d", rec.Code)
	}
}

func TestHandleAddWaterGet(t *testing.T) {
	h := testHandler(t)

	req := httptest.NewRequest("GET", "/water/add", nil)
	rec := httptest.NewRecorder()

	h.HandleAddWaterGet(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d", rec.Code)
	}
}

func TestHandleAddWaterPost(t *testing.T) {
	h := testHandler(t)

	form := url.Values{}
	form.Set("year", "2024")
	form.Set("volumeWater", "100")
	form.Set("volumeWastewater", "50")
	form.Set("volumeRainwater", "10")
	form.Set("costsWater", "20")
	form.Set("costsWastewater", "10")
	form.Set("costsRainwater", "5")
	form.Set("payments", "50")
	form.Set("fixedPrice", "15")

	req := httptest.NewRequest("POST", "/water/add", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.HandleAddWaterPost(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("got %d", rec.Code)
	}

	water, err := h.store.GetWater()
	if err != nil {
		t.Fatal(err)
	}

	if len(water) != 1 {
		t.Fatal("water missing")
	}

	if water[0].VolumeWater != 100 {
		t.Fatalf("got %f", water[0].VolumeWater)
	}
}

func TestHandleEditWater(t *testing.T) {
	h := testHandler(t)

	h.store.AddWater(models.WaterInput{Year: 2024})
	water, _ := h.store.GetWater()

	req := httptest.NewRequest("GET", "/water/edit", nil)
	req.SetPathValue("id", strconv.FormatInt(water[0].Id, 10))
	rec := httptest.NewRecorder()

	h.HandleEditWater(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d", rec.Code)
	}
}

func TestHandleEditWaterInvalidID(t *testing.T) {
	h := testHandler(t)

	req := httptest.NewRequest("GET", "/water/edit", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleEditWater(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got %d", rec.Code)
	}
}

func TestSaveWater(t *testing.T) {
	h := testHandler(t)

	h.store.AddWater(models.WaterInput{
		Year:        2024,
		VolumeWater: 100,
		CostsWater:  200,
	})

	waters, _ := h.store.GetWater()

	form := url.Values{}
	form.Set("year", "2025")
	form.Set("volumeWater", "300")
	form.Set("volumeWastewater", "250")
	form.Set("volumeRainwater", "10")
	form.Set("costsWater", "400")
	form.Set("costsWastewater", "100")
	form.Set("costsRainwater", "20")
	form.Set("payments", "500")
	form.Set("fixedPrice", "50")

	req := httptest.NewRequest("POST", "/water/save", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", strconv.FormatInt(waters[0].Id, 10))
	rec := httptest.NewRecorder()

	h.HandleSaveWater(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("got %d", rec.Code)
	}

	updated, _ := h.store.GetWater()
	if updated[0].Year != 2025 {
		t.Fatal("water not updated")
	}
	if updated[0].VolumeWater != 300 {
		t.Fatal("wrong water volume")
	}
}

func TestSaveWaterInvalidID(t *testing.T) {
	h := testHandler(t)

	form := url.Values{}
	form.Set("year", "2024")
	form.Set("volumeWater", "100")
	form.Set("volumeWastewater", "100")
	form.Set("volumeRainwater", "0")
	form.Set("costsWater", "200")
	form.Set("costsWastewater", "100")
	form.Set("costsRainwater", "0")
	form.Set("payments", "300")
	form.Set("fixedPrice", "0")

	req := httptest.NewRequest("POST", "/water/save", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleSaveWater(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got %d", rec.Code)
	}
}

func TestHandleDeleteWater(t *testing.T) {
	h := testHandler(t)

	h.store.AddWater(models.WaterInput{Year: 2024})
	water, _ := h.store.GetWater()

	req := httptest.NewRequest("GET", "/water/delete", nil)
	req.SetPathValue("id", strconv.FormatInt(water[0].Id, 10))
	rec := httptest.NewRecorder()

	h.HandleDeleteWater(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("got %d", rec.Code)
	}

	result, _ := h.store.GetWater()
	if len(result) != 0 {
		t.Fatal("water not deleted")
	}
}

func TestDeleteWaterInvalidID(t *testing.T) {
	h := testHandler(t)

	req := httptest.NewRequest("GET", "/water/delete", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.HandleDeleteWater(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got %d", rec.Code)
	}
}

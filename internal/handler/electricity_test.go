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

func TestElectricityView(t *testing.T) {

	h := testHandler(t)

	req :=
		httptest.NewRequest(
			"GET",
			"/electricity",
			nil,
		)

	rec :=
		httptest.NewRecorder()

	h.HandleElectricityView(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}

func TestAddElectricityGet(t *testing.T) {

	h := testHandler(t)

	req :=
		httptest.NewRequest(
			"GET",
			"/electricity/add",
			nil,
		)

	rec :=
		httptest.NewRecorder()

	h.HandleAddElectricityGet(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}

func TestAddElectricityPost(t *testing.T) {

	h := testHandler(t)

	form := url.Values{}

	form.Set(
		"timeFrom",
		"2024-01-01",
	)

	form.Set(
		"timeTo",
		"2024-12-31",
	)

	form.Set(
		"consumption",
		"3500",
	)

	form.Set(
		"costs",
		"1200",
	)

	form.Set(
		"payments",
		"1000",
	)

	form.Set(
		"retailer",
		"TestPower",
	)

	req :=
		httptest.NewRequest(
			"POST",
			"/electricity/add",
			strings.NewReader(
				form.Encode(),
			),
		)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	rec :=
		httptest.NewRecorder()

	h.HandleAddElectricityPost(
		rec,
		req,
	)

	if rec.Code != http.StatusFound {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}

	entries, err :=
		h.store.GetElectricities()

	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Fatal(
			"electricity missing",
		)
	}

	if entries[0].Consumption != 3500 {
		t.Fatal(
			"wrong consumption",
		)
	}
}

func TestEditElectricity(t *testing.T) {

	h := testHandler(t)

	err :=
		h.store.AddElectricity(
			models.ElectricityInput{
				TimeFrom: "2024",
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	entries, _ :=
		h.store.GetElectricities()

	req :=
		httptest.NewRequest(
			"GET",
			"/electricity/edit",
			nil,
		)

	req.SetPathValue(
		"id",
		strconv.FormatInt(
			entries[0].Id,
			10,
		),
	)

	rec :=
		httptest.NewRecorder()

	h.HandleEditElectricity(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}

func TestEditElectricityInvalidID(t *testing.T) {

	h := testHandler(t)

	req :=
		httptest.NewRequest(
			"GET",
			"/electricity/edit",
			nil,
		)

	req.SetPathValue(
		"id",
		"abc",
	)

	rec :=
		httptest.NewRecorder()

	h.HandleEditElectricity(
		rec,
		req,
	)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}

func TestSaveElectricity(t *testing.T) {

	h := testHandler(t)

	h.store.AddElectricity(
		models.ElectricityInput{
			TimeFrom:    "2024",
			Consumption: 100,
		},
	)

	entries, _ :=
		h.store.GetElectricities()

	form := url.Values{}

	form.Set(
		"timeFrom",
		"2025",
	)

	form.Set(
		"timeTo",
		"2025",
	)

	form.Set(
		"consumption",
		"200",
	)

	form.Set(
		"costs",
		"500",
	)

	form.Set(
		"payments",
		"300",
	)

	req :=
		httptest.NewRequest(
			"POST",
			"/electricity/save",
			strings.NewReader(
				form.Encode(),
			),
		)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	req.SetPathValue(
		"id",
		strconv.FormatInt(
			entries[0].Id,
			10,
		),
	)

	rec :=
		httptest.NewRecorder()

	h.HandleSaveElectricity(
		rec,
		req,
	)

	if rec.Code != http.StatusFound {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}

	updated, _ :=
		h.store.GetElectricities()

	if updated[0].Consumption != 200 {
		t.Fatal(
			"electricity not updated",
		)
	}
}

func TestSaveElectricityInvalidID(t *testing.T) {

	h := testHandler(t)

	form := url.Values{}

	form.Set(
		"timeFrom",
		"2025",
	)

	form.Set(
		"timeTo",
		"2025",
	)

	form.Set(
		"consumption",
		"100",
	)

	form.Set(
		"costs",
		"50",
	)

	form.Set(
		"payments",
		"20",
	)

	req :=
		httptest.NewRequest(
			"POST",
			"/electricity/save",
			strings.NewReader(
				form.Encode(),
			),
		)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	req.SetPathValue(
		"id",
		"abc",
	)

	rec :=
		httptest.NewRecorder()

	h.HandleSaveElectricity(
		rec,
		req,
	)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}

func TestDeleteElectricity(t *testing.T) {

	h := testHandler(t)

	h.store.AddElectricity(
		models.ElectricityInput{
			TimeFrom: "2024",
		},
	)

	entries, _ :=
		h.store.GetElectricities()

	req :=
		httptest.NewRequest(
			"GET",
			"/electricity/delete",
			nil,
		)

	req.SetPathValue(
		"id",
		strconv.FormatInt(
			entries[0].Id,
			10,
		),
	)

	rec :=
		httptest.NewRecorder()

	h.HandleDeleteElectricity(
		rec,
		req,
	)

	if rec.Code != http.StatusFound {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}

	result, _ :=
		h.store.GetElectricities()

	if len(result) != 0 {
		t.Fatal(
			"electricity not deleted",
		)
	}
}

func TestDeleteElectricityInvalidID(t *testing.T) {

	h := testHandler(t)

	req :=
		httptest.NewRequest(
			"GET",
			"/electricity/delete",
			nil,
		)

	req.SetPathValue(
		"id",
		"abc",
	)

	rec :=
		httptest.NewRecorder()

	h.HandleDeleteElectricity(
		rec,
		req,
	)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"got %d",
			rec.Code,
		)
	}
}

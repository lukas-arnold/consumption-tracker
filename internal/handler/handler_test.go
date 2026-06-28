package handler

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/storage"
)

func testHandler(t *testing.T) *Handler {

	t.Helper()

	store := storage.New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	return New(store)
}

func TestHandleView(t *testing.T) {

	h := testHandler(t)

	req :=
		httptest.NewRequest(
			"GET",
			"/",
			nil,
		)

	rec :=
		httptest.NewRecorder()

	h.HandleView(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}
}

func TestTemplateFuncs(t *testing.T) {

	funcs :=
		getTemplateFuncs()

	if funcs == nil {
		t.Fatal(
			"expected funcs",
		)
	}

	for name, fn := range funcs {

		if fn == nil {
			t.Fatalf(
				"template func %s nil",
				name,
			)
		}
	}
}

func TestTemplateTranslationFunc(t *testing.T) {

	funcs := getTemplateFuncs()

	tFunc :=
		funcs["T"].(func(string) string)

	result := tFunc("something")

	if result == "" {
		t.Fatal(
			"expected translation result",
		)
	}
}

func TestFormatFloat(t *testing.T) {

	result :=
		formatFloat(
			12.345,
			2,
		)

	if result != "12,35" {
		t.Fatalf(
			"got %s",
			result,
		)
	}
}

func TestFormatInt(t *testing.T) {

	result :=
		formatInt(
			12.7,
		)

	if result != "13" {
		t.Fatalf(
			"got %s",
			result,
		)
	}
}

func TestFormatEuro(t *testing.T) {

	if formatEuro(0) != "-" {
		t.Fatal(
			"expected dash",
		)
	}

	if formatEuro(12.5) != "12,50 €" {
		t.Fatal(
			"wrong euro format",
		)
	}
}

func TestFormatPrice(t *testing.T) {

	if formatPrice(
		10,
		0,
		"kWh",
	) != "-" {
		t.Fatal(
			"expected dash",
		)
	}

	result :=
		formatPrice(
			10,
			2,
			"kWh",
		)

	if result != "5,000 €/kWh" {
		t.Fatalf(
			"got %s",
			result,
		)
	}
}

func TestCalculateMonthlyPeriod(t *testing.T) {

	result :=
		calculateMonthlyPeriod(
			120,
			60,
		)

	if result != 60.875 {
		t.Fatalf(
			"got %f",
			result,
		)
	}
}

func TestFormatMonthly(t *testing.T) {

	if formatMonthly(
		0,
		"",
		"",
	) != "-" {
		t.Fatal(
			"expected dash",
		)
	}

	// fallback: no dates
	if formatMonthly(
		120,
		"",
		"",
	) != "10,00 €" {
		t.Fatal(
			"wrong monthly",
		)
	}
}

func TestFormatDifference(t *testing.T) {

	if formatDifference(
		10,
		10,
	) != "-" {
		t.Fatal(
			"expected dash",
		)
	}

	if formatDifference(
		100,
		40,
	) != "60,00 €" {
		t.Fatal(
			"wrong difference",
		)
	}
}

func TestFormatTotalCosts(t *testing.T) {

	if formatTotalCosts(
		0,
		0,
		0,
		0,
	) != "-" {
		t.Fatal(
			"expected dash",
		)
	}

	result :=
		formatTotalCosts(
			10,
			20,
			5,
			15,
		)

	if result != "50,00 €" {
		t.Fatalf(
			"got %s",
			result,
		)
	}
}

func TestAdd(t *testing.T) {

	result :=
		add(
			1,
			2,
			3,
		)

	if result != 6 {
		t.Fatalf(
			"got %f",
			result,
		)
	}
}

func TestRenderTemplate(t *testing.T) {

	h := testHandler(t)

	rec :=
		httptest.NewRecorder()

	h.renderTemplate(
		rec,
		"templates/electricity/add.html",
		nil,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}
}

func TestRenderTemplateWithData(t *testing.T) {

	h := testHandler(t)

	rec :=
		httptest.NewRecorder()

	h.renderTemplate(
		rec,
		"templates/electricity/index.html",
		models.Electricity{},
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}
}

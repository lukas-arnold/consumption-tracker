package handler

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/lukas-arnold/consumption-tracker/internal/storage"
)

func TestRoutes(t *testing.T) {

	store := storage.New(
		filepath.Join(
			t.TempDir(),
			"test.json",
		),
	)

	h := New(store)

	mux := http.NewServeMux()

	RegisterRoutes(
		mux,
		h,
	)

	tests := []struct {
		method string
		path   string
	}{
		{
			method: "GET",
			path:   "/",
		},

		{
			method: "GET",
			path:   "/web",
		},
		{
			method: "GET",
			path:   "/service-worker.js",
		},

		{
			method: "GET",
			path:   "/electricity",
		},
		{
			method: "GET",
			path:   "/electricity/add",
		},
		{
			method: "POST",
			path:   "/electricity/add",
		},
		{
			method: "GET",
			path:   "/electricity/edit/1",
		},
		{
			method: "POST",
			path:   "/electricity/save/1",
		},
		{
			method: "GET",
			path:   "/electricity/delete/1",
		},

		{
			method: "GET",
			path:   "/oil",
		},
		{
			method: "GET",
			path:   "/oil/add",
		},
		{
			method: "POST",
			path:   "/oil/add",
		},
		{
			method: "GET",
			path:   "/oil/edit/1",
		},
		{
			method: "POST",
			path:   "/oil/save/1",
		},
		{
			method: "GET",
			path:   "/oil/delete/1",
		},

		{
			method: "GET",
			path:   "/oil/fill-level/add",
		},
		{
			method: "POST",
			path:   "/oil/fill-level/add",
		},
		{
			method: "GET",
			path:   "/oil/fill-level/edit/1",
		},
		{
			method: "POST",
			path:   "/oil/fill-level/save/1",
		},
		{
			method: "GET",
			path:   "/oil/fill-level/delete/1",
		},

		{
			method: "GET",
			path:   "/water",
		},
		{
			method: "GET",
			path:   "/water/add",
		},
		{
			method: "POST",
			path:   "/water/add",
		},
		{
			method: "GET",
			path:   "/water/edit/1",
		},
		{
			method: "POST",
			path:   "/water/save/1",
		},
		{
			method: "GET",
			path:   "/water/delete/1",
		},
	}

	for _, test := range tests {

		req := httptest.NewRequest(
			test.method,
			test.path,
			nil,
		)

		rec := httptest.NewRecorder()

		mux.ServeHTTP(
			rec,
			req,
		)

		if rec.Code == http.StatusNotFound {
			t.Fatalf(
				"route missing: %s %s",
				test.method,
				test.path,
			)
		}
	}
}

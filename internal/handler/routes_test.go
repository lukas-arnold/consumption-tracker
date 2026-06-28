package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	h := testHandler(t)
	mux := http.NewServeMux()
	RegisterRoutes(mux, h)

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/"},
		{"GET", "/web/"},
		{"GET", "/service-worker.js"},

		{"GET", "/electricity"},
		{"GET", "/electricity/add"},
		{"POST", "/electricity/add"},
		{"GET", "/electricity/edit/1"},
		{"POST", "/electricity/save/1"},
		{"GET", "/electricity/delete/1"},

		{"GET", "/oil"},
		{"GET", "/oil/add"},
		{"POST", "/oil/add"},
		{"GET", "/oil/edit/1"},
		{"POST", "/oil/save/1"},
		{"GET", "/oil/delete/1"},

		{"GET", "/oil/fill-level/add"},
		{"POST", "/oil/fill-level/add"},
		{"GET", "/oil/fill-level/edit/1"},
		{"POST", "/oil/fill-level/save/1"},
		{"GET", "/oil/fill-level/delete/1"},

		{"GET", "/water"},
		{"GET", "/water/add"},
		{"POST", "/water/add"},
		{"GET", "/water/edit/1"},
		{"POST", "/water/save/1"},
		{"GET", "/water/delete/1"},
	}

	for _, tc := range tests {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code == http.StatusNotFound {
				t.Errorf("route %s %s returned 404", tc.method, tc.path)
			}
		})
	}
}

package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthAndTransportContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		name, method, path string
		check              ReadinessCheck
		status             int
		code               string
	}{
		{"live_without_database", "GET", "/health/live", nil, 200, ""},
		{"unconfigured_database", "GET", "/health/ready", nil, 503, "NOT_READY"},
		{"unavailable_database", "GET", "/health/ready", func(context.Context) error { return errors.New("private connection details") }, 503, "NOT_READY"},
		{"ready", "GET", "/health/ready", func(ctx context.Context) error {
			if _, ok := ctx.Deadline(); !ok {
				t.Error("readiness check lacks deadline")
			}
			return nil
		}, 200, ""},
		{"unknown_business_route", "GET", "/api/v1/admin/products", nil, 404, "NOT_FOUND"},
		{"invalid_health_method", "POST", "/health/live", nil, 405, "METHOD_NOT_ALLOWED"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			router := NewRouter(slog.New(slog.NewJSONHandler(io.Discard, nil)), tc.check)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
			if w.Code != tc.status {
				t.Fatalf("status=%d, want=%d; body=%s", w.Code, tc.status, w.Body)
			}
			if w.Header().Get("X-Request-ID") == "" {
				t.Fatal("missing request ID")
			}
			if tc.code != "" {
				var response ErrorResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				if response.Error.Code != tc.code || response.Error.RequestID != w.Header().Get("X-Request-ID") {
					t.Fatalf("invalid error contract: %+v", response)
				}
				if response.Error.Message == "private connection details" {
					t.Fatal("leaked dependency details")
				}
			}
		})
	}
}

func TestPanicReturnsGenericJSONError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(slog.New(slog.NewJSONHandler(io.Discard, nil)), nil)
	router.GET("/panic", func(c *gin.Context) { panic("secret") })
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/panic", nil))
	var response ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if w.Code != 500 || response.Error.Code != "INTERNAL_ERROR" {
		t.Fatalf("unexpected panic response: %d %s", w.Code, w.Body)
	}
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {
	tests := []struct {
		name           string
		platformID     string
		childProfileID string
		path           string
		expectedStatus int
		expectNext     bool
	}{
		{
			name:           "success - both headers present",
			platformID:     "vk",
			childProfileID: "550e8400-e29b-41d4-a716-446655440000",
			path:           "/test",
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name:           "success - only platform_id for onboarding",
			platformID:     "telegram",
			childProfileID: "",
			path:           "/onboarding/start",
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name:           "missing platform_id",
			platformID:     "",
			childProfileID: "",
			path:           "/test",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "missing child_profile_id for regular path",
			platformID:     "vk",
			childProfileID: "",
			path:           "/test",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test handler
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true

				// Check context values if expected to be present and path requires auth
				if tt.expectNext && requiresAuth(tt.path) {
					platformID := r.Context().Value(ContextKeyPlatformID)
					if platformID != tt.platformID {
						t.Errorf("expected platformID %q in context, got %v", tt.platformID, platformID)
					}

					if tt.childProfileID != "" {
						childProfileID := r.Context().Value(ContextKeyChildProfileID)
						if childProfileID != tt.childProfileID {
							t.Errorf("expected childProfileID %q in context, got %v", tt.childProfileID, childProfileID)
						}
					}
				}

				w.WriteHeader(http.StatusOK)
			})

			// Wrap with Auth middleware
			handler := Auth(next)

			// Create request
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.platformID != "" {
				req.Header.Set("X-Platform-ID", tt.platformID)
			}
			if tt.childProfileID != "" {
				req.Header.Set("X-Child-Profile-ID", tt.childProfileID)
			}

			w := httptest.NewRecorder()

			// Execute
			handler.ServeHTTP(w, req)

			// Assert
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if nextCalled != tt.expectNext {
				t.Errorf("expected nextCalled=%v, got %v", tt.expectNext, nextCalled)
			}
		})
	}
}

func TestCORS(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CORS(next)

	tests := []struct {
		name           string
		method         string
		origin         string
		expectedStatus int
		checkHeaders   func(t *testing.T, h http.Header)
	}{
		{
			name:           "preflight request",
			method:         http.MethodOptions,
			origin:         "http://localhost:5173",
			expectedStatus: http.StatusNoContent,
			checkHeaders: func(t *testing.T, h http.Header) {
				if h.Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
					t.Errorf("expected Access-Control-Allow-Origin to be http://localhost:5173, got %s", h.Get("Access-Control-Allow-Origin"))
				}
				if h.Get("Access-Control-Allow-Methods") == "" {
					t.Error("expected Access-Control-Allow-Methods to be set")
				}
				if h.Get("Access-Control-Allow-Headers") == "" {
					t.Error("expected Access-Control-Allow-Headers to be set")
				}
			},
		},
		{
			name:           "regular request",
			method:         http.MethodGet,
			origin:         "http://localhost:5173",
			expectedStatus: http.StatusOK,
			checkHeaders: func(t *testing.T, h http.Header) {
				if h.Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
					t.Errorf("expected Access-Control-Allow-Origin to be http://localhost:5173, got %s", h.Get("Access-Control-Allow-Origin"))
				}
			},
		},
		{
			name:           "disallowed origin",
			method:         http.MethodGet,
			origin:         "http://evil.com",
			expectedStatus: http.StatusOK,
			checkHeaders: func(t *testing.T, h http.Header) {
				if h.Get("Access-Control-Allow-Origin") != "" {
					t.Error("expected Access-Control-Allow-Origin to be empty for disallowed origin")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkHeaders != nil {
				tt.checkHeaders(t, w.Header())
			}
		})
	}
}

func TestRecovery(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		expectedStatus int
		expectPanic    bool
	}{
		{
			name: "no panic",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
			expectPanic:    false,
		},
		{
			name: "panic recovered",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic("test panic")
			},
			expectedStatus: http.StatusInternalServerError,
			expectPanic:    true,
		},
		{
			name: "panic with error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic(http.ErrAbortHandler)
			},
			expectedStatus: http.StatusInternalServerError,
			expectPanic:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Recovery(tt.handler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			// Execute (should not panic due to Recovery)
			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestLogging(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	handler := Logging(next)

	req := httptest.NewRequest(http.MethodGet, "/test?foo=bar", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Logging middleware should not modify response
	if w.Body.String() != "success" {
		t.Errorf("expected body 'success', got %q", w.Body.String())
	}
}

func TestChain(t *testing.T) {
	// Track execution order
	var executionOrder []string

	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware1-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware1-after")
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware2-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware2-after")
		})
	}

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, "handler")
		w.WriteHeader(http.StatusOK)
	})

	// Chain middlewares
	handler := Chain(middleware1, middleware2)(finalHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify execution order: middleware1 → middleware2 → handler → middleware2 → middleware1
	expectedOrder := []string{
		"middleware1-before",
		"middleware2-before",
		"handler",
		"middleware2-after",
		"middleware1-after",
	}

	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("execution[%d]: expected %q, got %q", i, expected, executionOrder[i])
		}
	}
}

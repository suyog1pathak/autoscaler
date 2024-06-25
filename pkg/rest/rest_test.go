package rest

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockHTTPClient is a mock implementation of the http.Client
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestClient(t *testing.T) {
	tests := []struct {
		method         string
		url            string
		headers        map[string]string
		body           []byte
		timeout        time.Duration
		responseBody   string
		expectedStatus int
		expectedBody   []byte
	}{
		{
			method:         "GET",
			url:            "/test",
			headers:        map[string]string{"Content-Type": "application/json"},
			body:           nil,
			timeout:        time.Second * 10,
			responseBody:   `{"message":"success"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"message":"success"}`),
		},
		{
			method:         "POST",
			url:            "/test",
			headers:        map[string]string{"Content-Type": "application/json"},
			body:           []byte(`{"key":"value"}`),
			timeout:        time.Second * 10,
			responseBody:   `{"status":"created"}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   []byte(`{"status":"created"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.method+" "+test.url, func(t *testing.T) {
			// Setup a mock server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(test.expectedStatus)
				io.WriteString(w, test.responseBody)
			}))
			defer ts.Close()

			response, responseBody, err := Client(test.method, ts.URL+test.url, test.headers, test.body, test.timeout)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if response.StatusCode != test.expectedStatus {
				t.Errorf("expected status %d, got %d", test.expectedStatus, response.StatusCode)
			}

			if !bytes.Equal(responseBody, test.expectedBody) {
				t.Errorf("expected body %s, got %s", test.expectedBody, responseBody)
			}
		})
	}
}

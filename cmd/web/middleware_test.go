package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zaki-Zak/Snippet-Go-Box/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	reqrec := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("okay"))
	})

	commonHeaders(next).ServeHTTP(reqrec, r)

	result := reqrec.Result()

	// INFO: Content-Security-Policy header check
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, result.Header.Get("Content-Security-Policy"), expectedValue)

	// INFO: Referrer-Policy header check
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, result.Header.Get("Referrer-Policy"), expectedValue)

	// INFO: X-Content-Type-Options header check
	expectedValue = "nosniff"
	assert.Equal(t, result.Header.Get("X-Content-Type-Options"), expectedValue)

	// INFO: X-Frame-Options header check
	expectedValue = "deny"
	assert.Equal(t, result.Header.Get("X-Frame-Options"), expectedValue)

	// INFO: X-XSS-Protection header check
	expectedValue = "0"
	assert.Equal(t, result.Header.Get("X-XSS-Protection"), expectedValue)

	// INFO: Server header check
	expectedValue = "Go"
	assert.Equal(t, result.Header.Get("Server"), expectedValue)

	// INFO: next handler called check
	assert.Equal(t, result.StatusCode, http.StatusOK)

	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "okay")
}

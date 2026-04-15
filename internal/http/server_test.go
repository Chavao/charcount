package http_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chavao/charcount/internal/config"
	apphttp "github.com/Chavao/charcount/internal/http"
)

func TestNewServerServesIndex(t *testing.T) {
	t.Parallel()

	server := apphttp.NewServer(config.Config{
		Host: "127.0.0.1",
		Port: 5536,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	server.Handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode = %d, want %d", response.StatusCode, http.StatusOK)
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", contentType)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	bodyText := string(body)
	if !strings.Contains(bodyText, "charcount-app") {
		t.Fatalf("body missing app root: %q", bodyText)
	}

	if !strings.Contains(bodyText, "<textarea") {
		t.Fatalf("body missing textarea: %q", bodyText)
	}

	if !strings.Contains(bodyText, "Word density") {
		t.Fatalf("body missing density section: %q", bodyText)
	}

	if !strings.Contains(bodyText, `/assets/app.js`) {
		t.Fatalf("body missing app script reference: %q", bodyText)
	}
}

func TestNewServerServesEmbeddedAssets(t *testing.T) {
	t.Parallel()

	server := apphttp.NewServer(config.Config{
		Host: "127.0.0.1",
		Port: 5536,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)

	server.Handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode = %d, want %d", response.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if !strings.Contains(string(body), "analyzeText") {
		t.Fatalf("body missing script content: %q", string(body))
	}
}

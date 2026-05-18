package http_test

import (
	"bytes"
	"encoding/json"
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

	if !strings.Contains(string(body), "/api/analyze") {
		t.Fatalf("body missing script content: %q", string(body))
	}
}

func TestAnalyzeTextEndpoint(t *testing.T) {
	t.Parallel()

	server := apphttp.NewServer(config.Config{
		Host: "127.0.0.1",
		Port: 5536,
	})

	requestBody := bytes.NewBufferString(`{"text":"Hello world! Hello."}`)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/analyze", requestBody)
	request.Header.Set("Content-Type", "application/json")

	server.Handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode = %d, want %d", response.StatusCode, http.StatusOK)
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", contentType)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	var payload struct {
		Characters  int `json:"characters"`
		Words       int `json:"words"`
		Sentences   int `json:"sentences"`
		Paragraphs  int `json:"paragraphs"`
		Spaces      int `json:"spaces"`
		DensityRows []struct {
			Word    string `json:"word"`
			Count   int    `json:"count"`
			Density string `json:"density"`
		} `json:"densityRows"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if payload.Characters != len("Hello world! Hello.") {
		t.Fatalf("Characters = %d, want %d", payload.Characters, len("Hello world! Hello."))
	}

	if payload.Words != 3 {
		t.Fatalf("Words = %d, want 3", payload.Words)
	}

	if payload.Sentences != 2 {
		t.Fatalf("Sentences = %d, want 2", payload.Sentences)
	}

	if payload.Paragraphs != 1 {
		t.Fatalf("Paragraphs = %d, want 1", payload.Paragraphs)
	}

	if payload.Spaces != 2 {
		t.Fatalf("Spaces = %d, want 2", payload.Spaces)
	}

	if len(payload.DensityRows) != 2 {
		t.Fatalf("DensityRows length = %d, want 2", len(payload.DensityRows))
	}

	if payload.DensityRows[0].Word != "hello" || payload.DensityRows[0].Count != 2 || payload.DensityRows[0].Density != "66.67%" {
		t.Fatalf("DensityRows[0] = %+v, want hello x2 66.67%%", payload.DensityRows[0])
	}

	if payload.DensityRows[1].Word != "world" || payload.DensityRows[1].Count != 1 || payload.DensityRows[1].Density != "33.33%" {
		t.Fatalf("DensityRows[1] = %+v, want world x1 33.33%%", payload.DensityRows[1])
	}
}

func TestAnalyzeTextEndpointRejectsInvalidJSON(t *testing.T) {
	t.Parallel()

	server := apphttp.NewServer(config.Config{
		Host: "127.0.0.1",
		Port: 5536,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/analyze", strings.NewReader(`{"text":"ok","unknown":1}`))
	request.Header.Set("Content-Type", "application/json")

	server.Handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode = %d, want %d", response.StatusCode, http.StatusBadRequest)
	}
}

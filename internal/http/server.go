package http

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Chavao/charcount/internal/analyzer"
	"github.com/Chavao/charcount/internal/config"
	"github.com/Chavao/charcount/internal/web"
)

func NewServer(cfg config.Config) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", index)
	mux.HandleFunc("POST /api/analyze", analyzeText)
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServerFS(web.AssetFS())))

	return &http.Server{
		Addr:              cfg.Address(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(web.IndexHTML)
}

type analyzeRequest struct {
	Text string `json:"text"`
}

func analyzeText(w http.ResponseWriter, r *http.Request) {
	requestBody := http.MaxBytesReader(w, r.Body, 1<<20)
	defer requestBody.Close()

	decoder := json.NewDecoder(requestBody)
	decoder.DisallowUnknownFields()

	var payload analyzeRequest
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var extraData struct{}
	if err := decoder.Decode(&extraData); err != io.EOF {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	analysis := analyzer.Analyze(payload.Text)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(analysis); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

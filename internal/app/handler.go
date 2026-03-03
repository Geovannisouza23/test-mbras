package app

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"backend-challenge-092025/internal/domain"
)

// Handler handles HTTP requests.
type Handler struct {
	Processor *domain.Processor
}

func NewHandler(p *domain.Processor) *Handler {
	return &Handler{Processor: p}
}

func (h *Handler) AnalyzeFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	var req domain.AnalyzeFeedRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()
	resp, code, msg := h.Processor.AnalyzeFeed(req, now)
	if code == 400 {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if code == 422 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(msg))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

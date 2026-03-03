package app

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"backend-challenge-092025/internal/domain"
	errpkg "backend-challenge-092025/pkg/errors"
	validatorpkg "backend-challenge-092025/pkg/validator"
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
		errpkg.BadRequest(w, []byte(`{"error": "Method Not Allowed"}`))
		return
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		errpkg.BadRequest(w, []byte(`{"error": "Unsupported Media Type"}`))
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errpkg.BadRequest(w, errpkg.RespJSONDecodeFailure)
		return
	}
	var req domain.AnalyzeFeedRequest
	if err := json.Unmarshal(body, &req); err != nil {
		errpkg.BadRequest(w, errpkg.RespJSONDecodeFailure)
		return
	}
	if err := validatorpkg.ValidateStruct(req); err != nil {
		respBody, _ := json.Marshal(validatorpkg.ToErrResponse(err))
		errpkg.ValidationErrors(w, respBody)
		return
	}
	now := time.Now().UTC()
	resp, code, msg := h.Processor.AnalyzeFeed(req, now)
	if code == 400 {
		errpkg.BadRequest(w, []byte(msg))
		return
	}
	if code == 422 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(msg))
		return
	}
	errpkg.WriteJSON(w, http.StatusOK, resp)
}

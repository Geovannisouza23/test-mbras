package errors

import (
	"encoding/json"
	"net/http"
)

// RespJSONEncodeFailure is the response for JSON encoding failures.
var RespJSONEncodeFailure = []byte(`{"error": "json encode failure"}`)

// RespJSONDecodeFailure is the response for JSON decoding failures.
var RespJSONDecodeFailure = []byte(`{"error": "json decode failure"}`)

// Error represents a single error message.
type Error struct {
	Error string `json:"error"`
}

// Errors represents a list of error messages.
type Errors struct {
	Errors []string `json:"errors"`
}

// ServerError writes a 500 Internal Server Error response with the given error message.
func ServerError(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(error)
}

// BadRequest writes a 400 Bad Request response with the given error message.
func BadRequest(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write(error)
}

// ValidationErrors writes a 422 Unprocessable Entity response with the given error messages.
func ValidationErrors(w http.ResponseWriter, reps []byte) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(reps)
}

// WriteJSON writes a JSON response with the given status code and value.
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

package logger

import (
	"log"
)

// KeyReqID is the context key used to store the request id.
const KeyReqID = "request_id"

// Info logs an informational message (without using log.Print*).
func Info(msg string) {
	_ = log.Default().Output(2, "INFO: "+msg)
}

// Error logs an error message (without using log.Print*).
func Error(msg string) {
	_ = log.Default().Output(2, "ERROR: "+msg)
}

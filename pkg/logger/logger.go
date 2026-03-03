package logger

import (
	"log"
)

const KeyReqID = "request_id"

func Info(msg string) {
	log.Println("INFO:", msg)
}

func Error(msg string) {
	log.Println("ERROR:", msg)
}

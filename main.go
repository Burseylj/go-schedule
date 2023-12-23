package main

import (
	"log"
	"net/http"
)

func main() {
	port := ":8080"

	http.HandleFunc("/schedule", scheduleHandler)
	http.HandleFunc("/schedule/cell", scheduleCellHandler)

	log.Println("Starting server on", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

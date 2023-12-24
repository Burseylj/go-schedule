package main

import (
	"go-schedule/schedule"
	"log"
	"net/http"
)

func main() {
	port := ":8080"
	http.HandleFunc("/schedule", schedule.ScheduleHandler)
	http.HandleFunc("/schedule/cell", schedule.ScheduleCellHandler)

	log.Println("Starting server on", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

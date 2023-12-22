package main

import (
	"cloud.google.com/go/civil"
	"log"
	"net/http"
)

//init test data
var employees = []Employee{
	{1, "Alice Smith", "TeamA"},
	{2, "Bob Johnson", "TeamB"},
	{3, "Charlie Brown", "TeamA"},
	{4, "Diana Davis", "TeamB"},
	{5, "Edward Wilson", "TeamA"},
	{6, "Frank Thompson", "TeamA"},
	{7, "Grace Martinez", "TeamB"},
	{8, "Henry Anderson", "TeamA"},
}

var dates = []civil.Date{
	{Year: 2023, Month: 1, Day: 1},
	{Year: 2023, Month: 1, Day: 2},
	{Year: 2023, Month: 1, Day: 3},
	{Year: 2023, Month: 1, Day: 4},
	{Year: 2023, Month: 1, Day: 5},
}

var employeeSchedule = Schedule{
	"1:2023-01-01": "Day",
	"4:2023-01-05": "Eve",
}

func main() {
	http.HandleFunc("/schedule", scheduleHandler)
	http.HandleFunc("/schedule/cell", scheduleCellHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

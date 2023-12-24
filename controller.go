package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"cloud.google.com/go/civil"
)


var scheduleTemplate *template.Template

//init test data, will be db service eventually
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

func httpError(w http.ResponseWriter, error string, code int) {
    log.Println(error)
    http.Error(w, error, code)
}


func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpError(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	group := groupByTeam(employees)

	err := schedule(group, dates, employeeSchedule).Render(r.Context(), w)
	if err != nil {
		httpError(w, "Error executing template: %v", http.StatusBadRequest)
	}
}

func groupByTeam(employees []Employee) map[string][]Employee {
	groups := make(map[string][]Employee)
	for _, emp := range employees {
		groups[emp.Team] = append(groups[emp.Team], emp)
	}
	return groups
}

func parseEmpID(w http.ResponseWriter, r *http.Request) (int, error) {
	param := r.URL.Query().Get("empID")
	parsed, err := strconv.Atoi(param)
	if err != nil {
		httpError(w, "Error parsing "+param, http.StatusBadRequest)
		return 0, err
	}
	return parsed, nil
}

func parseDate(w http.ResponseWriter, r *http.Request) (civil.Date, error) {
	param := r.URL.Query().Get("date")
	parsed, err := civil.ParseDate(param)
	if err != nil {
		httpError(w, "Error parsing "+param, http.StatusBadRequest)
		return civil.Date{}, err
	}
	return parsed, nil
}

func scheduleCellHandler(w http.ResponseWriter, r *http.Request) {
	empID, err := parseEmpID(w, r)
	if err != nil {
		return
	}
	date, err := parseDate(w, r)
	if err != nil {
		return
	}

	switch r.Method {
	case "PUT":
		cellContent := r.FormValue("event")
		employeeSchedule.Update(empID, date, cellContent)
		renderCellContent(w, r, empID, date, cellContent)

	case "DELETE":
		employeeSchedule.Delete(empID, date)
		renderCellContent(w, r, empID, date, "")

	default:
		httpError(w, "Only PUT and DELETE methods are allowed", http.StatusMethodNotAllowed)	
	}
}

func renderCellContent(w http.ResponseWriter, r *http.Request, empID int, date civil.Date, content string) {
	err := cellContents(empID, date, content).Render(r.Context(), w)
	if err != nil {
		httpError(w, "Error rendering cell", http.StatusInternalServerError)
		log.Printf("Error rendering cell: %v", err)
	}
}
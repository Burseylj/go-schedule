package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"cloud.google.com/go/civil"
	"github.com/Masterminds/sprig"
)


var scheduleTemplate *template.Template

func init() {
	var err error
	scheduleTemplate, err = template.New("schedule.html").Funcs(sprig.FuncMap()).ParseFiles("schedule.html")
	if err != nil {
		log.Fatalf("Error loading template: %v", err)
	}
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

	if err := renderSchedulePage(w); err != nil {
		httpError(w, "Error executing template: %v", http.StatusBadRequest)
		log.Printf("Error executing template: %v", err)
	}
}

func renderSchedulePage(w http.ResponseWriter) error {
	data := struct {
		Employees []Employee
		Dates     []civil.Date
		Schedule  Schedule
	}{
		Employees: employees,
		Dates:     dates,
		Schedule:  employeeSchedule,
	}

	return scheduleTemplate.Execute(w, data)
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
		renderCellContent(w, empID, date, cellContent)

	case "DELETE":
		employeeSchedule.Delete(empID, date)
		renderCellContent(w, empID, date, "")

	default:
		httpError(w, "Only PUT and DELETE methods are allowed", http.StatusMethodNotAllowed)	
	}
}

func renderCellContent(w http.ResponseWriter, empID int, date civil.Date, content string) {
	data := map[string]interface{}{
		"Content": content,
		"EmpID":   empID,
		"Date":    date,
	}

	err := scheduleTemplate.ExecuteTemplate(w, "cellContent", data)
	if err != nil {
		httpError(w, "Error rendering cell", http.StatusInternalServerError)
		log.Printf("Error rendering cell: %v", err)
	}
}
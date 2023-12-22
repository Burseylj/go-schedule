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

func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	renderSchedulePage(w)
}

func renderSchedulePage(w http.ResponseWriter) {
	data := struct {
		Employees []Employee
		Dates     []civil.Date
		Schedule  Schedule
	}{
		Employees: employees,
		Dates:     dates,
		Schedule:  employeeSchedule,
	}

	if err := scheduleTemplate.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func parseEmpID(w http.ResponseWriter, r *http.Request) (int, error) {
	param := r.URL.Query().Get("empID")
	parsed, err := strconv.Atoi(param)
	if err != nil {
		http.Error(w, "Error parsing "+param, http.StatusBadRequest)
		log.Println("Error parsing", param, err)
		return 0, err
	}
	return parsed, nil
}

func parseDate(w http.ResponseWriter, r *http.Request) (civil.Date, error) {
	param := r.URL.Query().Get("date")
	parsed, err := civil.ParseDate(param)
	if err != nil {
		http.Error(w, "Error parsing "+param, http.StatusBadRequest)
		log.Println("Error parsing", param, err)
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
		log.Printf("Attempting to add cell: %v", err)
		cellContent := r.FormValue("event")
		employeeSchedule.Update(empID, date, cellContent)
		renderCellContent(w, empID, date, cellContent)

	case "DELETE":
		log.Printf("Attempting to delete cell: %v", err)
		employeeSchedule.Delete(empID, date)
		renderCellContent(w, empID, date, "")

	default:
		http.Error(w, "Only PUT and DELETE methods are allowed", http.StatusMethodNotAllowed)
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
		http.Error(w, "Error rendering cell", http.StatusInternalServerError)
		log.Printf("Error rendering cell: %v", err)
	}
}
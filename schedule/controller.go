package schedule

import (
	"cloud.google.com/go/civil"
	"go-schedule/model"
	"log"
	"net/http"
	"strconv"
)

// init test data, will be db service eventually
var employees = []model.Employee{
	{ID: 1, Name: "Alice Smith", Team: "TeamA"},
	{ID: 2, Name: "Bob Johnson", Team: "TeamB"},
	{ID: 3, Name: "Charlie Brown", Team: "TeamA"},
	{ID: 4, Name: "Diana Davis", Team: "TeamB"},
	{ID: 5, Name: "Edward Wilson", Team: "TeamA"},
	{ID: 6, Name: "Frank Thompson", Team: "TeamA"},
	{ID: 7, Name: "Grace Martinez", Team: "TeamB"},
	{ID: 8, Name: "Henry Anderson", Team: "TeamA"},
}

var dates = []civil.Date{
	{Year: 2023, Month: 1, Day: 1},
	{Year: 2023, Month: 1, Day: 2},
	{Year: 2023, Month: 1, Day: 3},
	{Year: 2023, Month: 1, Day: 4},
	{Year: 2023, Month: 1, Day: 5},
}

var employeeSchedule = model.NewSchedule()

func init() {
	employeeSchedule.Set(1, dates[0], "Vacation")
	employeeSchedule.Set(2, dates[1], "Day")
}

func httpError(w http.ResponseWriter, error string, code int) {
	log.Println(error)
	http.Error(w, error, code)
}

func ScheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpError(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	group := groupByTeam(employees)

	err := scheduleTemplate(group, dates, employeeSchedule).Render(r.Context(), w)
	if err != nil {
		httpError(w, "Error executing template: %v", http.StatusBadRequest)
	}
}

func groupByTeam(employees []model.Employee) map[string][]model.Employee {
	groups := make(map[string][]model.Employee)
	for _, emp := range employees {
		groups[emp.Team] = append(groups[emp.Team], emp)
	}
	return groups
}

func parseEmpID(w http.ResponseWriter, r *http.Request) (int, error) {
	param := r.URL.Query().Get("empID")
	parsed, err := strconv.Atoi(param)
	if err != nil {
		httpError(w, "Error parsing empID"+param, http.StatusBadRequest)
		return 0, err
	}
	return parsed, nil
}

func parseDate(w http.ResponseWriter, r *http.Request) (civil.Date, error) {
	param := r.URL.Query().Get("date")
	parsed, err := civil.ParseDate(param)
	if err != nil {
		httpError(w, "Error parsing Date"+param, http.StatusBadRequest)
		return civil.Date{}, err
	}
	return parsed, nil
}

func ScheduleCellHandler(w http.ResponseWriter, r *http.Request) {
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
		employeeSchedule.Set(empID, date, cellContent)
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

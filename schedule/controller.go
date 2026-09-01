package schedule

import (
	"cloud.google.com/go/civil"
	"log"
	"net/http"
	"strconv"
)

type Employee struct {
	ID   int
	Name string
	Team string
}

// init test data, will be db service eventually
var employees = []Employee{
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


type Store interface {
	Get(empID int, date civil.Date) string
	Set(empID int, date civil.Date, event string)
	Delete(empID int, date civil.Date)
}

var store = NewMapStore()

func init() {
	store.Set(1, dates[0], "Vacation")
	store.Set(2, dates[1], "Day")
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

	err := scheduleTemplate(group, dates, store).Render(r.Context(), w)
	if err != nil {
		httpError(w, "Error executing template\n" + err.Error(), http.StatusInternalServerError)
		return
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
		httpError(w, "Error parsing empID " + param, http.StatusBadRequest)
		return 0, err
	}
	return parsed, nil
}

func parseDate(w http.ResponseWriter, r *http.Request) (civil.Date, error) {
	param := r.URL.Query().Get("date")
	parsed, err := civil.ParseDate(param)
	if err != nil {
		httpError(w, "Error parsing Date " + param, http.StatusBadRequest)
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
		store.Set(empID, date, cellContent)
		renderCellContent(w, r, empID, date, cellContent)

	case "DELETE":
		store.Delete(empID, date)
		renderCellContent(w, r, empID, date, "")

	default:
		httpError(w, "Only PUT and DELETE methods are allowed", http.StatusMethodNotAllowed)
	}
}

func renderCellContent(w http.ResponseWriter, r *http.Request, empID int, date civil.Date, content string) {
	err := cellContents(empID, date, content).Render(r.Context(), w)
	if err != nil {
		httpError(w, "Error rendering cell", http.StatusInternalServerError)
	}
}

package model

import (
	"cloud.google.com/go/civil"
	"fmt"
)

type Employee struct {
	ID   int
	Name string
	Team string
}

type schedule map[string]string

type Schedule interface {
	Set(empID int, date civil.Date, event string)
	Get(empID int, date civil.Date) string
	Delete(empID int, date civil.Date)
}

func NewSchedule() Schedule {
	return &schedule{}
}

func (s *schedule) Set(empID int, date civil.Date, event string) {
	key := fmt.Sprintf("%d:%s", empID, date)
	(*s)[key] = event
}

func (s *schedule) Get(empID int, date civil.Date) string {
	key := fmt.Sprintf("%d:%s", empID, date)
	return (*s)[key]
}

func (s *schedule) Delete(empID int, date civil.Date) {
	key := fmt.Sprintf("%d:%s", empID, date)
	delete(*s, key)
}

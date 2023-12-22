package main

import (
	"fmt"
	"cloud.google.com/go/civil"
)

type Employee struct {
	ID   int
	Name string
	Team string
}

type Schedule map[string]string

func (s Schedule) Update(empID int, date civil.Date, event string) {
	key := fmt.Sprintf("%d:%s", empID, date)
	s[key] = event
}

func (s Schedule) Get(empID int, date civil.Date) string {
	key := fmt.Sprintf("%d:%s", empID, date)
	return s[key]
}

func (s Schedule) Delete(empID int, date civil.Date) {
	key := fmt.Sprintf("%d:%s", empID, date)
	delete(s, key)
}
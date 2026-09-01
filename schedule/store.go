package schedule

import (
	"sync"

	"cloud.google.com/go/civil"
)

type key struct {
	EmpID int
	Date  civil.Date
}

type MapStore struct {
	mu   sync.RWMutex
	data map[key]string
}

func NewMapStore() *MapStore {
	return &MapStore{data: make(map[key]string)}
}

func (s *MapStore) Get(empID int, date civil.Date) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[key{EmpID: empID, Date: date}]
}

func (s *MapStore) Set(empID int, date civil.Date, event string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key{EmpID: empID, Date: date}] = event
}

func (s *MapStore) Delete(empID int, date civil.Date) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key{EmpID: empID, Date: date})
}
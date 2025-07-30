package common

import (
	"slices"
	"sync"
)

type StringList struct {
	List []string
	Lock *sync.RWMutex
}

// NewStringList
func NewStringList() *StringList {
	return &StringList{List: make([]string, 0), Lock: &sync.RWMutex{}}
}

// Add
func (s *StringList) Add(str string) bool {
	if s.Contains(str) {
		return true
	}
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.List = append(s.List, str)
	return true
}

// Remove
func (s *StringList) Remove(str string) bool {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	for i, v := range s.List {
		if v == str {
			s.List = slices.Delete(s.List, i, i+1)
			return true
		}
	}
	return false
}

// Contains
func (s *StringList) Contains(str string) bool {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return slices.Contains(s.List, str)
}

// Len
func (s *StringList) Len() int {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return len(s.List)
}

// Clear
func (s *StringList) Clear() {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.List = make([]string, 0)
}

// get
func (s *StringList) Get(index int) string {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	if index < 0 || index >= len(s.List) {
		if len(s.List) > 0 {
			return s.List[0]
		}
		return ""
	}
	return s.List[index]
}

// Range
func (s *StringList) Range(f func(index int, str string) bool) bool {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	for i, v := range s.List {
		if !f(i, v) {
			return false
		}
	}
	return true
}

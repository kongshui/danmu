package common

import "sync"

type StringMap struct {
	Map  map[string]string
	Lock *sync.RWMutex
}

// newStringmap
func NewStringMap() *StringMap {
	return &StringMap{Map: make(map[string]string), Lock: &sync.RWMutex{}}
}

// add
func (s *StringMap) Add(key, value string) {
	if s.Contains(key) {
		return
	}
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Map[key] = value
}

// remove
func (s *StringMap) Remove(key string) {
	if !s.Contains(key) {
		return
	}
	s.Lock.Lock()
	defer s.Lock.Unlock()
	delete(s.Map, key)
}

// get
func (s *StringMap) Get(key string) string {
	if !s.Contains(key) {
		return ""
	}
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	value := s.Map[key]
	return value
}

// get all
func (s *StringMap) GetAll() map[string]string {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return s.Map
}

// contains key
func (s *StringMap) Contains(key string) bool {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	_, ok := s.Map[key]
	return ok
}

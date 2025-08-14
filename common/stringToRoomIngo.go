package common

import "sync"

type (
	StringToRoomInfoMap struct {
		Map  map[string]RoomInfo
		Lock *sync.RWMutex
	}
	RoomInfo struct {
		RoomId string //roomid
		UserId string // openId
	}
)

// NewStringToRoomInfoMap
func NewStringToRoomInfoMap() *StringToRoomInfoMap {
	return &StringToRoomInfoMap{Map: make(map[string]RoomInfo), Lock: &sync.RWMutex{}}
}

// Add
func (s *StringToRoomInfoMap) Add(key string, value RoomInfo) bool {
	if s.Contains(key) {
		return true
	}
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Map[key] = value
	return true
}

// get one Key
func (s *StringToRoomInfoMap) Get(key string) RoomInfo {
	if s.Contains(key) {
		return s.Map[key]
	}
	return RoomInfo{}
}

// GetAll Key

func (s *StringToRoomInfoMap) GetAll() map[string]RoomInfo {
	return s.Map
}

// Remove
func (s *StringToRoomInfoMap) Remove(key string) bool {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if _, ok := s.Map[key]; ok {
		delete(s.Map, key)
		return true
	}
	return false
}

func (s *StringToRoomInfoMap) Contains(key string) bool {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	_, ok := s.Map[key]
	return ok
}

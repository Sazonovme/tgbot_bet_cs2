package safemap

import "sync"

type SafeMap struct {
	mu           sync.RWMutex
	lastMsgIDMap map[int64]int // map[ChatID]messageID
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		lastMsgIDMap: make(map[int64]int),
	}
}

func (s *SafeMap) Get(chatID int64) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.lastMsgIDMap[chatID]
	return val, ok
}

func (s *SafeMap) Set(chatID int64, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastMsgIDMap[chatID] = value
}

func (s *SafeMap) Delete(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.lastMsgIDMap, chatID)
}

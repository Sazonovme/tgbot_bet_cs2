package safemap

import (
	"sync"
	"time"
)

type UserSession struct {
	LastMsgIDs []int // список ID сообщений
	SendAt     time.Time
	State      string // текущее состояние (main_menu, matches_list и т.д.)
}

type SafeMap struct {
	mu       sync.RWMutex
	sessions map[int64]UserSession // map[ChatID]session
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		sessions: make(map[int64]UserSession),
	}
}

func (s *SafeMap) Get(chatID int64) ([]int, time.Time, string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.sessions[chatID]
	return val.LastMsgIDs, val.SendAt, val.State, ok
}

func (s *SafeMap) Set(chatID int64, lastMsgIDs []int, state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[chatID] = UserSession{
		LastMsgIDs: lastMsgIDs,
		State:      state,
		SendAt:     time.Now(),
	}
}

func (s *SafeMap) Delete(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, chatID)
}

func (s *SafeMap) ChangeLastMessages(chatID int64, lastMsgIDs []int, state string) {
	s.Delete(chatID)
	s.Set(chatID, lastMsgIDs, state)
}

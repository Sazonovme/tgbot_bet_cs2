package usersessions

import (
	"sync"
	"time"
)

type UserSession struct {
	LastMsgIDs []int // список ID сообщений
	SendAt     time.Time
	State      string // текущее состояние (main_menu, matches_list и т.д.)
}

type UserSessionMap struct {
	mu       sync.RWMutex
	sessions map[int64]UserSession // map[ChatID]session
}

func NewUserSessionMap() *UserSessionMap {
	return &UserSessionMap{
		sessions: make(map[int64]UserSession),
	}
}

func (s *UserSessionMap) Get(chatID int64) ([]int, time.Time, string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.sessions[chatID]
	return val.LastMsgIDs, val.SendAt, val.State, ok
}

func (s *UserSessionMap) Set(chatID int64, lastMsgIDs []int, state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[chatID] = UserSession{
		LastMsgIDs: lastMsgIDs,
		State:      state,
		SendAt:     time.Now(),
	}
}

func (s *UserSessionMap) Delete(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, chatID)
}

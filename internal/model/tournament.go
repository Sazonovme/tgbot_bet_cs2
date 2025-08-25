package model

import (
	"database/sql"
	"slices"
	"time"
)

var AdminChatIDs []int64

func IsAdmin(chat_id int64) bool {
	return slices.Contains(AdminChatIDs, chat_id)
}

type UserData struct {
	Chat_id int64
	Text    string
	IsAdmin bool
}

type ConfirmPrediction struct {
	MatchName string
	Match_id  int
	Tag       string
	Bet       string
	BetText   string
	TextMsg   string
	Confirmed bool
}

type Match struct {
	Id    int
	Name  string
	Date  time.Time
	Team1 string
	Team2 string
}

type Result struct {
	Match_id int
	Result   string
}

type UserPrediction struct {
	Username   string
	Chat_id    int64
	Match_id   int
	Match_name string
	Prediction string
	DateMatch  time.Time
	Result     sql.NullString
}

type TournamentFinishTable struct {
	Username        string
	Match_name      string
	Match_date      time.Time
	User_prediction string
	Match_result    string
	Score           int
}

type ScoreFinishTable map[string]int

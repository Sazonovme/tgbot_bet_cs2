package model

import "time"

type User struct {
	Chat_id      int64
	Username     string
	First_name   string
	Last_name    string
	TextMsg      string
	CallbackData string
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
	Match_Name string
	Username   string
	Match_id   uint
	Prediction string
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

package model

import "time"

type User struct {
	Chat_id    int64
	Username   string
	First_name string
	Last_name  string
	TextMsg    string
}

type Event struct {
	Name  string
	Date  time.Time
	Team1 string
	Team2 string
}

type UserPrediction struct {
	NameEvent  string
	UserName   string
	Id_event   uint
	Prediction string
}

type EventFinishTable struct {
	Username        string
	Name_match      string
	Date_match      time.Time
	User_prediction string
	Result_match    string
	Score           int
}

type ScoreFinishTable map[string]int

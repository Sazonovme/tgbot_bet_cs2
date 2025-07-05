package model

import "time"

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

type FinishTable struct {
	Username   string
	Name       string
	Prediction string
	Result     string
}

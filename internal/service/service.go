package service

import (
	"RushBananaBet/internal/model"
	"RushBananaBet/pkg/logger"
	"context"
	"strings"
)

type Service struct {
	Repository Repository
}

type Repository interface {
	CreateEvent(ctx context.Context, event model.Event) error
	AddResultToEvent(ctx context.Context, result string) error
	GetEventFinishTable(ctx context.Context) ([]model.EventFinishTable, error)
	GetUserPredictions(ctx context.Context, username string) ([]model.UserPrediction, error)
	AddUserPrediction(ctx context.Context, prediction *model.UserPrediction) error
	AddNewUser(ctx context.Context, user *model.User) error
	DeactivateUser(ctx context.Context, chat_id int64) error
}

func NewService(repo Repository) *Service {
	return &Service{
		Repository: repo,
	}
}

// GENERAL

func (s *Service) Start(ctx context.Context, user *model.User) error {
	return s.Repository.AddNewUser(ctx, user)
}

func (s *Service) Stop(ctx context.Context, chat_id int64) error {
	return s.Repository.DeactivateUser(ctx, chat_id)
}

// ADMIN

func (s *Service) CreateEvent(ctx context.Context, event model.Event) error {
	return s.Repository.CreateEvent(ctx, event)
}

func (s *Service) AddResultToEvent(ctx context.Context, result string) error {
	return s.Repository.AddResultToEvent(ctx, result)
}

func (s *Service) GetEventFinishTable(ctx context.Context) ([]model.EventFinishTable, model.ScoreFinishTable, error) {
	eventFinishTable, err := s.Repository.GetEventFinishTable(ctx)
	if err != nil {
		logger.Error("Error get event finish table from repo", "service-GetEventFinishTable()", err)
		return nil, nil, err
	}

	mScore := map[string]int{}
	curScore := 0
	for key, elem := range eventFinishTable {

		// Calc score
		if strings.Contains(elem.User_prediction, "t") {
			prediction_team := strings.Replace(elem.User_prediction, "t", "", -1)
			firstChar := string([]rune(elem.Result_match)[0])
			if (firstChar == "2" && prediction_team == "1") || (firstChar != "2" && prediction_team == "2") {
				curScore = 1
			}
		} else {
			if elem.User_prediction == elem.Result_match {
				curScore = 2
			}
		}

		// Add score
		eventFinishTable[key].Score = curScore
		if curScore != 0 {
			val, ok := mScore[elem.Username]
			if !ok {
				mScore[elem.Username] = curScore
			} else {
				mScore[elem.Username] = val + curScore
			}
			curScore = 0
		}
	}

	return eventFinishTable, mScore, nil
}

// USER

func (s *Service) GetUserPredictions(ctx context.Context, username string) ([]model.UserPrediction, error) {
	return s.Repository.GetUserPredictions(ctx, username)
}

func (s *Service) AddUserPrediction(ctx context.Context, prediction *model.UserPrediction) error {
	return s.Repository.AddUserPrediction(ctx, prediction)
}

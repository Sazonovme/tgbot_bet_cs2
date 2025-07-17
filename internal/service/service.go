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
	CreateTournament(ctx context.Context, name_tournament string) error
	CreateMatch(ctx context.Context, match *model.Match) error
	AddMatchResult(ctx context.Context, result string, match_id int) error
	GetTournamentFinishTable(ctx context.Context) (*[]model.TournamentFinishTable, error)
	GetUserPredictions(ctx context.Context, username string) (*[]model.UserPrediction, error)
	AddUserPrediction(ctx context.Context, prediction *model.UserPrediction) error
	AddNewUser(ctx context.Context, user *model.User) error
	DeactivateUser(ctx context.Context, chat_id int64) error
}

func NewService(repo Repository) *Service {
	return &Service{
		Repository: repo,
	}
}

// ADMIN

func (s *Service) CreateTournament(ctx context.Context, name_tournament string) error {
	return s.Repository.CreateTournament(ctx, name_tournament)
}

func (s *Service) CreateMatch(ctx context.Context, match *model.Match) error {
	return s.Repository.CreateMatch(ctx, match)
}

func (s *Service) AddMatchResult(ctx context.Context, result string, match_id int) error {
	return s.Repository.AddMatchResult(ctx, result, match_id)
}

func (s *Service) GetTournamentFinishTable(ctx context.Context) (*[]model.TournamentFinishTable, *model.ScoreFinishTable, error) {
	tournamentFinishTablePointer, err := s.Repository.GetTournamentFinishTable(ctx)
	if err != nil {
		logger.Error("Error get event finish table from repo", "service-GetEventFinishTable()", err)
		return nil, nil, err
	}

	tournamentFinishTable := *tournamentFinishTablePointer
	mScore := model.ScoreFinishTable{}
	curScore := 0
	for key, elem := range tournamentFinishTable {

		// Calc score
		if strings.Contains(elem.User_prediction, "t") {
			prediction_team := strings.Replace(elem.User_prediction, "t", "", -1)
			firstChar := string([]rune(elem.Match_result)[0])
			if (firstChar == "2" && prediction_team == "1") || (firstChar != "2" && prediction_team == "2") {
				curScore = 1
			}
		} else {
			if elem.User_prediction == elem.Match_result {
				curScore = 2
			}
		}

		// Add score
		tournamentFinishTable[key].Score = curScore
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

	return &tournamentFinishTable, &mScore, nil
}

// USER

func (s *Service) GetUserPredictions(ctx context.Context, username string) (*[]model.UserPrediction, error) {
	return s.Repository.GetUserPredictions(ctx, username)
}

func (s *Service) AddUserPrediction(ctx context.Context, prediction *model.UserPrediction) error {
	return s.Repository.AddUserPrediction(ctx, prediction)
}

// GENERAL

func (s *Service) AddNewUser(ctx context.Context, user *model.User) error {
	return s.Repository.AddNewUser(ctx, user)
}

func (s *Service) DeactivateUser(ctx context.Context, chat_id int64) error {
	return s.Repository.DeactivateUser(ctx, chat_id)
}

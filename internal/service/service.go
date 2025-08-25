package service

import (
	"RushBananaBet/internal/logger"
	"RushBananaBet/internal/model"
	"context"
	"strings"
)

type Service struct {
	Repository Repository
}

type Repository interface {

	//GENERAL
	AddNewUser(ctx context.Context, chat_id int64, user_id int64, username string) (isExist bool, err error)

	// ADMIN
	CreateTournament(ctx context.Context, name_tournament string) (added bool, err error)
	CreateMatches(ctx context.Context, matches []model.Match) error
	AddMatchResults(ctx context.Context, results []model.Result) error
	GetTournamentFinishTable(ctx context.Context) ([]model.TournamentFinishTable, error)
	GetActiveMatchesID(ctx context.Context) ([]model.Match, error)

	// USER
	GetActiveMatches(ctx context.Context) ([]model.Match, error)
	GetUserPredictions(ctx context.Context, chat_id int64) ([]model.UserPrediction, error)
	AddUpdateUserPrediction(ctx context.Context, chat_id int64, match_id int, prediction string) (inserted bool, err error)
}

func NewService(repo Repository) *Service {
	return &Service{
		Repository: repo,
	}
}

// GENERAL

func (s *Service) AddNewUser(ctx context.Context, chat_id int64, user_id int64, username string) (isExist bool, err error) {
	return s.Repository.AddNewUser(ctx, chat_id, user_id, username)
}

// ADMIN

func (s *Service) CreateTournament(ctx context.Context, name_tournament string) (added bool, err error) {
	return s.Repository.CreateTournament(ctx, name_tournament)
}

func (s *Service) CreateMatches(ctx context.Context, matches []model.Match) error {
	return s.Repository.CreateMatches(ctx, matches)
}

func (s *Service) AddMatchResults(ctx context.Context, results []model.Result) error {
	return s.Repository.AddMatchResults(ctx, results)
}

func (s *Service) GetTournamentFinishTable(ctx context.Context) ([]model.TournamentFinishTable, model.ScoreFinishTable, error) {

	tournamentFinishTable, err := s.Repository.GetTournamentFinishTable(ctx)
	if err != nil {
		logger.Error("Error get finish table", "service-GetTournamentFinishTable()", err)
		return nil, nil, err
	}

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

	logger.Debug("Success final table calculation", "service-GetTournamentFinishTable()", nil)
	return tournamentFinishTable, mScore, nil
}

func (s *Service) GetActiveMatchesID(ctx context.Context) ([]model.Match, error) {
	return s.Repository.GetActiveMatchesID(ctx)
}

// USER

func (s *Service) GetActiveMatches(ctx context.Context) ([]model.Match, error) {
	return s.Repository.GetActiveMatches(ctx)
}

func (s *Service) GetUserPredictions(ctx context.Context, chat_id int64) ([]model.UserPrediction, error) {
	return s.Repository.GetUserPredictions(ctx, chat_id)
}

func (s *Service) AddUpdateUserPrediction(ctx context.Context, chat_id int64, match_id int, prediction string) (inserted bool, err error) {
	return s.Repository.AddUpdateUserPrediction(ctx, chat_id, match_id, prediction)
}

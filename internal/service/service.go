package service

import (
	"RushBananaBet/internal/logger"
	"RushBananaBet/internal/model"
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Repository Repository
}

type Repository interface {
	CreateTournament(ctx context.Context, name_tournament string) error
	CreateMatch(ctx context.Context, matches *[]model.Match) error
	AddMatchResult(ctx context.Context, results *[]model.Result) error
	GetTournamentFinishTable(ctx context.Context) (*[]model.TournamentFinishTable, error)
	GetMatchesIDs(ctx context.Context) (*[]model.Match, error)
	GetActiveMatches(ctx context.Context) (*[]model.Match, error)
	GetUserPredictions(ctx context.Context, username string) (*[]model.UserPrediction, error)
	AddUserPrediction(ctx context.Context, prediction *model.UserPrediction) error
	AddNewUser(ctx context.Context, user *model.User) (err error, isExist bool)
	DeactivateUser(ctx context.Context, chat_id int64) error
}

func NewService(repo Repository) *Service {
	return &Service{
		Repository: repo,
	}
}

// ADMIN

func (s *Service) CreateTournament(ctx context.Context, userData *model.User) error {
	args := strings.Split(userData.TextMsg, " ")
	if len(args) < 2 {
		return errors.New("arr message is clear")
	}
	return s.Repository.CreateTournament(ctx, args[1])
}

func (s *Service) CreateMatch(ctx context.Context, userData *model.User) error {
	matches := []model.Match{}
	args := strings.Split(userData.TextMsg, " ")
	matches_string_array := strings.Split(args[1], "#")

	for _, val := range matches_string_array {
		match_arr := strings.Split(val, "_")
		teams_arr := strings.Split(match_arr[0], "vs")

		date_match, err := time.Parse("02.01.2006 03:04", match_arr[1])
		if err != nil {
			logger.Error("Err parse time in create match", "service-CreateMatch()", err)
			continue
		}

		matches = append(matches, model.Match{
			Name:  match_arr[0],
			Date:  date_match,
			Team1: teams_arr[0],
			Team2: teams_arr[1],
		})
	}

	if len(matches) < 1 {
		return errors.New("error parse matches")
	}

	return s.Repository.CreateMatch(ctx, &matches)
}

func (s *Service) AddMatchResult(ctx context.Context, userData *model.User) error {
	results := []model.Result{}
	args := strings.Split(userData.TextMsg, " ")
	result_string_array := strings.Split(args[1], "#")

	for _, val := range result_string_array {
		result_arr := strings.Split(val, "_")
		match_id, err := strconv.Atoi(result_arr[0])

		if err != nil {
			logger.Error("Err convert string match_id to int", "service-AddMatchResult()", err)
			continue
		}

		results = append(results, model.Result{
			Match_id: match_id,
			Result:   result_arr[1],
		})
	}

	return s.Repository.AddMatchResult(ctx, &results)
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

func (s *Service) GetMatchesIDs(ctx context.Context) (*[]model.Match, error) {
	return s.Repository.GetMatchesIDs(ctx)
}

// USER

func (s *Service) GetActiveMatches(ctx context.Context) (*[]model.Match, error) {
	return s.Repository.GetActiveMatches(ctx)
}

func (s *Service) GetUserPredictions(ctx context.Context, username string) (*[]model.UserPrediction, error) {
	return s.Repository.GetUserPredictions(ctx, username)
}

func (s *Service) AddUserPrediction(ctx context.Context, userData *model.User) error {

	// make_prediction_[matchID]_[bet]_[y/n]
	args := strings.Split(userData.CallbackData, "_")
	match_id, err := strconv.Atoi(args[2])
	if err != nil {
		logger.Error("Err to convert string to int", "service - AddUserPrediction()", err)
		return err
	}

	prediction := model.UserPrediction{
		Match_id:   uint(match_id),
		Prediction: args[3],
	}

	return s.Repository.AddUserPrediction(ctx, &prediction)
}

// GENERAL

func (s *Service) AddNewUser(ctx context.Context, user *model.User) (err error, isExist bool) {
	return s.Repository.AddNewUser(ctx, user)
}

func (s *Service) DeactivateUser(ctx context.Context, chat_id int64) error {
	return s.Repository.DeactivateUser(ctx, chat_id)
}

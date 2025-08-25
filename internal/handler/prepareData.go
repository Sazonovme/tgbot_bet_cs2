package handler

import (
	"RushBananaBet/internal/logger"
	"RushBananaBet/internal/model"
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidData = errors.New("invalid format data")
	ErrParseTime   = errors.New("invalid format date")
)

func prepareCreateTournamentData(data string) (tournamentName string, err error) {
	str := strings.TrimSpace(data)
	if len(str) == 0 {
		return "", ErrInvalidData
	}
	return data, nil
}

func prepareCreateMatchesData(data string) (matches []model.Match, err error) {

	//[t1_t2]-[16.08.2025 15:00]#...
	result := []model.Match{}
	matches_string_array := strings.Split(data, "#")

	for _, val := range matches_string_array {

		match_arr := strings.Split(val, "-")
		if len(match_arr) != 2 {
			return nil, ErrInvalidData
		}

		teams_arr := strings.Split(match_arr[0], "_")
		if len(teams_arr) != 2 {
			return nil, ErrInvalidData
		}

		date_match, err := time.Parse("02.01.2006 03:04", match_arr[1])
		if err != nil {
			logger.Error("Err parse time in prepare for create match", "handler-prepareCreateMatchesData()", err)
			return nil, ErrParseTime
		}

		result = append(result, model.Match{
			Name:  teams_arr[0] + " vs " + teams_arr[1],
			Date:  date_match,
			Team1: teams_arr[0],
			Team2: teams_arr[1],
		})
	}
	return result, nil
}

func prepareAddMatchResultsData(data string) (results []model.Result, err error) {

	//[matchID]_[result]#...
	resultsArr := []model.Result{}
	result_string_array := strings.Split(data, "#")

	for _, val := range result_string_array {

		arr := strings.Split(val, "_")
		if len(arr) != 2 {
			return nil, ErrInvalidData
		}

		match_id, err := strconv.Atoi(arr[0])
		if err != nil {
			logger.Error("Err convert string match_id to int", "handler-prepareAddMatchResults()", err)
			return nil, ErrInvalidData
		}

		resultsArr = append(resultsArr, model.Result{
			Match_id: match_id,
			Result:   arr[1],
		})
	}
	return resultsArr, nil
}

func prepareConfirmPredictionData(data string) (confirm_predictions model.ConfirmPrediction, err error) {

	// confirm_prediction_[matchName]_[matchID]_[bet]
	// change_prediction_[matchName]_[matchID]_[bet]

	arr := strings.Split(data, "_")

	// Match_id
	match_id, err := strconv.Atoi(arr[3])
	if err != nil {
		return model.ConfirmPrediction{}, ErrInvalidData
	}

	// Readable bet for user
	betTxt := getReadableBet(arr[4])

	// Result message
	textMessage := ""
	if arr[0] == "confirm" {
		textMessage = "Матч: " + arr[2] + "\n" + "Ваша ставка: " + betTxt + "\n" + "Подтвердить ставку?"
	} else {
		textMessage = "Изменение ставки\n" + "Матч: " + arr[2] + "\n" + "Новая ставка: " + betTxt + "\n" + "Подтвердить?"
	}

	return model.ConfirmPrediction{
		MatchName: arr[2],
		Match_id:  match_id,
		Tag:       "End" + arr[0],
		BetText:   betTxt,
		Bet:       arr[4],
		TextMsg:   textMessage,
	}, nil
}

func prepareProcessingConfirmPredictionData(data string) (confirm_predictions model.ConfirmPrediction, err error) {

	// Endchange_prediction_[matchID]_[bet]_[y/n]
	arr := strings.Split(data, "_")

	// match_id
	match_id, err := strconv.Atoi(arr[2])
	if err != nil {
		return model.ConfirmPrediction{}, ErrInvalidData
	}

	// Confirmed
	confirmed := false
	if arr[4] == "y" {
		confirmed = true
	}

	return model.ConfirmPrediction{
		Match_id:  match_id,
		Tag:       arr[0],
		Bet:       arr[3],
		Confirmed: confirmed,
	}, nil

}

func getReadableBet(data string) string {
	if data == "1" || data == "2" {
		return "Победа команды " + data
	} else {
		return "Точный счет " + data
	}
}

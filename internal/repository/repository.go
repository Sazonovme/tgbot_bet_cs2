package repository

import (
	"RushBananaBet/internal/logger"
	"RushBananaBet/internal/model"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type mainRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *mainRepository {
	return &mainRepository{
		db: db,
	}
}

// GENERAL

func (r *mainRepository) AddNewUser(ctx context.Context, chat_id int64, user_id int64, username string) (isExist bool, err error) {

	query := `INSERT INTO telegram_users (chat_id, user_id, username, is_active)
				VALUES (@chat_id, @user_id, @username, true)`
	args := pgx.NamedArgs{
		"chat_id":  chat_id,
		"user_id":  user_id,
		"username": username,
	}

	_, err = r.db.Exec(ctx, query, args)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				logger.Debug("User alredy exist", "repository-AddNewUser()", nil)
				return true, nil
			}
		} else {
			logger.Error("Error add new user in db", "repository-AddNewUser()", err)
			return false, err
		}
	}

	logger.Debug("Success add new user in db", "repository-AddNewUser()", nil)
	return false, nil
}

// ADMIN

func (r *mainRepository) CreateTournament(ctx context.Context, name_tournament string) (added bool, err error) {

	query := `WITH inserted AS (
			  INSERT INTO tournaments (name, is_active)
			  SELECT $1, true
			  WHERE NOT EXISTS (
				  SELECT 1
				  FROM tournaments
				  WHERE is_active = true
			  )
			  RETURNING id
			  )
			  SELECT id FROM inserted`

	rows, err := r.db.Query(ctx, query, name_tournament)
	if err != nil {
		logger.Error("Err r.db.Query()", "repository-CreateTournament()", err)
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		logger.Debug("Success create new tournament", "repository-CreateTournament()", nil)
		return true, nil
	} else {
		logger.Debug("Active tournament already exist", "repository-CreateTournament()", nil)
		return false, nil
	}

}

func (r *mainRepository) CreateMatches(ctx context.Context, matches []model.Match) error {

	valueStrings := make([]string, 0, len(matches))
	valueArgs := make([]any, 0, len(matches)*4)

	for i, match := range matches {
		n := i * 4
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, CAST($%d AS timestamptz))", n+1, n+2, n+3, n+4))
		valueArgs = append(valueArgs, match.Name, match.Team1, match.Team2, match.Date)
	}

	query := fmt.Sprintf(`
	WITH t AS (
		SELECT id FROM tournaments WHERE is_active = true LIMIT 1
	),
	match_data(name, team_1, team_2, date) AS (
		VALUES %s
	)
	INSERT INTO matches (tournament_id, name, team_1, team_2, date)
	SELECT t.id, m.name, m.team_1, m.team_2, m.date
	FROM t, match_data m`, strings.Join(valueStrings, ","))

	_, err := r.db.Exec(ctx, query, valueArgs...)
	if err != nil {
		logger.Error("Create matches in db error", "repository-CreateMatches()", err)
		return err
	}

	logger.Debug("Success create matches in db", "repository-CreateMatches()", nil)
	return nil
}

func (r *mainRepository) AddMatchResults(ctx context.Context, results []model.Result) error {

	valueStrings := make([]string, 0, len(results))
	valueArgs := make([]any, 0, len(results)*2)

	for i, r := range results {
		n := i * 2
		valueStrings = append(valueStrings, fmt.Sprintf("($%d::int, $%d::text)", n+1, n+2))
		valueArgs = append(valueArgs, r.Match_id, r.Result)
	}

	query := fmt.Sprintf(`
		UPDATE matches AS m
		SET result = v.result
		FROM (
			VALUES %s
		) AS v(id, result)
		WHERE m.id = v.id;
	`, strings.Join(valueStrings, ","))

	_, err := r.db.Exec(ctx, query, valueArgs...)
	if err != nil {
		logger.Error("Add match results in db error", "repository-AddMatchResults()", err)
		return err
	}

	logger.Debug("Success add match results in db", "repository-AddMatchResults()", nil)
	return nil
}

func (r *mainRepository) GetTournamentFinishTable(ctx context.Context) ([]model.TournamentFinishTable, error) {

	query := `SELECT
    			tu.username,
    			m.name AS match_name,
   				p.prediction,
    			m.result,
    			m.date
			FROM predictions p
			JOIN telegram_users tu ON p.chat_id = tu.chat_id
			JOIN matches m ON p.match_id = m.id
			JOIN tournaments t ON m.tournament_id = t.id
			WHERE tu.is_active = true
				AND t.is_active = true
			ORDER BY m.date ASC`

	sqlRows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.Error("Get finish table in db error", "repository-getTournamentFinishTable()", err)
		return nil, err
	}
	defer sqlRows.Close()

	finishTables := []model.TournamentFinishTable{}
	for sqlRows.Next() {
		finishTable := model.TournamentFinishTable{}
		err := sqlRows.Scan(&finishTable.Username, &finishTable.Match_name, &finishTable.User_prediction, &finishTable.Match_result, &finishTable.Match_date)
		if err != nil {
			logger.Error("Scan finish table in db error", "repository-getTournamentFinishTable()", err)
			return nil, err
		}
		finishTables = append(finishTables, finishTable)
	}

	logger.Debug("Success get finish table in db", "repository-getTournamentFinishTable()", nil)
	return finishTables, nil
}

func (r *mainRepository) GetActiveMatchesID(ctx context.Context) ([]model.Match, error) {

	query := `SELECT
    			id,
    			name,
   				date
			FROM matches
			WHERE result IS NULL
			ORDER BY date ASC`

	sqlRows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.Error("Get matches IDs in db error", "repository-GetActiveMatchesID()", err)
		return nil, err
	}
	defer sqlRows.Close()

	finishArr := []model.Match{}
	for sqlRows.Next() {
		elem := model.Match{}
		err := sqlRows.Scan(&elem.Id, &elem.Name, &elem.Date)
		if err != nil {
			logger.Error("Scan match IDs in db error", "repository-GetActiveMatchesID()", err)
			return nil, err
		}
		finishArr = append(finishArr, elem)
	}

	logger.Debug("Success get match IDs in db", "repository-GetActiveMatchesID()", nil)
	return finishArr, nil
}

// USER

func (r *mainRepository) GetActiveMatches(ctx context.Context) ([]model.Match, error) {

	query := `SELECT
				id,
				name,
				team_1,
				team_2,
				date			
			FROM matches
			WHERE
				date > @currentDate`
	args := pgx.NamedArgs{
		"currentDate": time.Now().Add(120 * time.Second),
	}

	sqlRows, err := r.db.Query(ctx, query, args)
	if err != nil {
		logger.Error("Get active matches in db error", "repository-GetActiveMatches()", err)
		return nil, err
	}
	defer sqlRows.Close()

	activeMatches := []model.Match{}
	for sqlRows.Next() {
		activeMatch := model.Match{}
		err := sqlRows.Scan(&activeMatch.Id, &activeMatch.Name, &activeMatch.Team1, &activeMatch.Team2, &activeMatch.Date)
		if err != nil {
			logger.Error("Scan active matches in db error", "repository-GetActiveMatches()", err)
			return nil, err
		}
		activeMatches = append(activeMatches, activeMatch)
	}

	logger.Debug("Success get active matches in db", "repository-GetActiveMatches()", nil)
	return activeMatches, nil
}

func (r *mainRepository) GetUserPredictions(ctx context.Context, chat_id int64) ([]model.UserPrediction, error) {

	query := `SELECT
				m.id AS match_id,
				m.name AS match_name,
				p.prediction,
				m.date,
				m.result			
			FROM predictions p
			JOIN matches m ON p.match_id = m.id
			WHERE
				p.chat_id = @chat_id
			ORDER BY
    			CASE 
       				WHEN m.result is NULL THEN 0          	  -- сначала матчи с результатом
        			WHEN m.date < NOW() THEN 1                -- затем прошедшие матчи без результата
       				ELSE 2                                    -- потом будущие матчи без результата
				END,
				m.date ASC`
	args := pgx.NamedArgs{
		"chat_id": chat_id,
	}

	sqlRows, err := r.db.Query(ctx, query, args)
	if err != nil {
		logger.Error("Get user predictions in db error", "repository-GetUserPredictions()", err)
		return nil, err
	}
	defer sqlRows.Close()

	predictions := []model.UserPrediction{}
	for sqlRows.Next() {
		prediction := model.UserPrediction{}
		err := sqlRows.Scan(&prediction.Match_id, &prediction.Match_name, &prediction.Prediction, &prediction.DateMatch, &prediction.Result)
		if err != nil {
			logger.Error("Scan prediction in db error", "repository-GetUserPredictions()", err)
			return nil, err
		}
		predictions = append(predictions, prediction)
	}

	logger.Debug("Success get user predictions in db", "repository-GetUserPredictions()", nil)
	return predictions, err
}

func (r *mainRepository) AddUpdateUserPrediction(ctx context.Context, chat_id int64, match_id int, prediction string) (inserted bool, err error) {

	query := `INSERT INTO predictions (chat_id, match_id, prediction)
				VALUES (
					$1,
					$2, 
					$3
				)
				ON CONFLICT (chat_id, match_id)
				DO UPDATE SET prediction = EXCLUDED.prediction
				RETURNING 
    			CASE WHEN xmax = 0 THEN 'inserted' ELSE 'updated' END AS action`

	var action string
	err = r.db.QueryRow(ctx, query, chat_id, match_id, prediction).Scan(&action)
	if err != nil {
		logger.Error("Add/update user predictions in db error", "repository-AddUpdateUserPrediction()", err)
		return false, err
	}

	if action == "inserted" {
		logger.Debug("Success add user prediction in db", "repository-AddUpdateUserPrediction()", nil)
		return true, nil
	} else {
		logger.Debug("Success update user prediction in db", "repository-AddUpdateUserPrediction()", nil)
		return false, nil
	}
}

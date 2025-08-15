package repository

import (
	"RushBananaBet/internal/logger"
	"RushBananaBet/internal/model"
	"context"
	"errors"
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

// ADMIN

func (r *mainRepository) CreateTournament(ctx context.Context, name_tournament string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		logger.Error("Error new transaction in create tournaments", "repository-CreateTournament()", err)
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `UPDATE tournaments SET is_active = false WHERE is_active = true`)
	if err != nil {
		logger.Error("Error in sql set is_active=false for tournament", "repository-CreateTournament()", err)
		return err
	}

	query := `INSERT INTO tournaments(name, is_active) VALUES (@name, true)`
	args := pgx.NamedArgs{
		"name": name_tournament,
	}
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		logger.Error("Error in sql insert new tournament", "repository-CreateTournament()", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("Error commit tx", "repository-CreateTournament()", err)
		return err
	}

	logger.Debug("Success create new tournament", "repository-CreateTournament()", nil)
	return nil
}

func (r *mainRepository) CreateMatch(ctx context.Context, matches *[]model.Match) error {

	if len(*matches) == 0 {
		return errors.New("match array is clear")
	}

	valueStrings := make([]string, 0, len(*matches))
	valueArgs := make([]interface{}, 0, len(*matches)*4)
	for i, match := range *matches {
		n := i * 5
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, %d)", n+1, n+2, n+3, n+4, n+5))
		valueArgs = append(valueArgs, match.Name, match.Team1, match.Team2, match.Date, "")
	}

	query := fmt.Sprintf(`
	WITH t AS (
		SELECT id FROM tournaments WHERE is_active = true LIMIT 1
	),
	match_data(name, team_1, team_2, date, result) AS (
		VALUES %s
	)
	INSERT INTO matches (tournament_id, name, team_1, team_2, date, result)
	SELECT t.id, m.name, m.team_1, m.team_2, m.date, m.result
	FROM t, match_data m`, strings.Join(valueStrings, ","))

	_, err := r.db.Exec(ctx, query, valueArgs...)
	if err != nil {
		logger.Error("Create match in db error", "repository-createEvent()", err)
		return err
	}
	logger.Debug("Success create match in db", "repository-createEvent()", nil)
	return nil
}

func (r *mainRepository) AddMatchResult(ctx context.Context, results *[]model.Result) error {

	if len(*results) == 0 {
		return errors.New("result array is clear")
	}

	valueStrings := make([]string, 0, len(*results))
	valueArgs := make([]interface{}, 0, len(*results)*2)

	for i, r := range *results {
		n := i * 2
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", n+1, n+2))
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
		logger.Error("Add result to match in db error", "repository-AddMatchResult()", err)
		return err
	}
	logger.Debug("Success add result match in db", "repository-AddMatchResult()", nil)
	return nil
}

func (r *mainRepository) GetTournamentFinishTable(ctx context.Context) (*[]model.TournamentFinishTable, error) {

	query := `SELECT
    			tu.username,
    			m.name AS match_name,
   				p.prediction,
    			m.result,
    			m.date
			FROM predictions p
			JOIN telegram_users tu ON p.username = tu.username
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

	logger.Debug("Success get sql rows for finish table in db", "repository-getTournamentFinishTable()", nil)

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
	return &finishTables, nil
}

func (r *mainRepository) GetMatchesIDs(ctx context.Context) (*[]model.Match, error) {

	query := `SELECT
    			id,
    			name,
   				date
			FROM matches
			WHERE result = ''
			ORDER BY date ASC`

	sqlRows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.Error("Get matches IDs in db error", "repository-GetMatchesIDs()", err)
		return nil, err
	}
	defer sqlRows.Close()

	logger.Debug("Success get sql rows match IDs in db", "repository-GetMatchesIDs()", nil)

	finishArr := []model.Match{}
	for sqlRows.Next() {
		elem := model.Match{}
		err := sqlRows.Scan(&elem.Id, &elem.Name, &elem.Date)
		if err != nil {
			logger.Error("Scan match IDs in db error", "repository-GetMatchesIDs()", err)
			return nil, err
		}
		finishArr = append(finishArr, elem)
	}

	logger.Debug("Success get match IDs in db", "repository-GetMatchesIDs()", nil)
	return &finishArr, nil
}

// USER

func (r *mainRepository) GetActiveMatches(ctx context.Context) (*[]model.Match, error) {

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
		"currentDate": time.Now().Add(60 * time.Second),
	}

	sqlRows, err := r.db.Query(ctx, query, args)
	if err != nil {
		logger.Error("Get active matches in db error", "repository-GetActiveMatches()", err)
		return nil, err
	}
	defer sqlRows.Close()

	logger.Debug("Success get active matches in db", "repository-GetActiveMatches()", nil)

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
	return &activeMatches, nil
}

func (r *mainRepository) GetUserPredictions(ctx context.Context, username string) (*[]model.UserPrediction, error) {

	query := `SELECT
				m.name AS match_name,
				p.prediction,
				m.date,
				m.result			
			FROM predictions p
			JOIN matches m ON p.match_id = m.id
			WHERE
				p.username = @username
			ORDER BY
    			CASE 
       				WHEN m.result <> '' THEN 0           -- сначала матчи с результатом
        			WHEN m.date < NOW() THEN 1                -- затем прошедшие матчи без результата
       				ELSE 2                                    -- потом будущие матчи без результата
				END,
				m.date ASC`
	args := pgx.NamedArgs{
		"username": username,
	}

	sqlRows, err := r.db.Query(ctx, query, args)
	if err != nil {
		logger.Error("Get user predictions in db error", "repository-getUserPredictions()", err)
		return nil, err
	}
	defer sqlRows.Close()

	logger.Debug("Success get sql rows for user predictions in db", "repository-getUserPredictions()", nil)

	predictions := []model.UserPrediction{}
	for sqlRows.Next() {
		prediction := model.UserPrediction{}
		err := sqlRows.Scan(&prediction.Match_Name, &prediction.Prediction, &prediction.DateMatch, &prediction.Result)
		if err != nil {
			logger.Error("Scan prediction in db error", "repository-getUserPredictions()", err)
			return nil, err
		}
		predictions = append(predictions, prediction)
	}

	logger.Debug("Success get user predictions in db", "repository-getUserPredictions()", nil)
	return &predictions, nil
}

func (r *mainRepository) AddUserPrediction(ctx context.Context, prediction *model.UserPrediction) error {

	query := `INSERT INTO predictions (username, match_id, prediction)
				VALUES (@username, @match_id, @prediction)
				ON CONFLICT (username, match_id)
				DO UPDATE SET prediction = EXCLUDED.prediction`
	args := pgx.NamedArgs{
		"username":   prediction.Username,
		"match_id":   prediction.Match_id,
		"prediction": prediction.Prediction,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		logger.Error("Add user prediction in db error", "repository-addUserPrediction()", err)
		return err
	}

	logger.Debug("Success add user prediction in db", "repository-addUserPrediction()", nil)
	return nil
}

func (r *mainRepository) AddNewUser(ctx context.Context, user *model.User) (err error, isExist bool) {
	query := `INSERT INTO telegram_users (chat_id, username, first_name, last_name, is_active)
				VALUES (@chat_id, @username, @first_name, @last_name, true)`
	args := pgx.NamedArgs{
		"chat_id":    user.Chat_id,
		"username":   user.Username,
		"first_name": user.First_name,
		"last_name":  user.Last_name,
	}

	_, err = r.db.Exec(ctx, query, args)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				return nil, true
			}
		} else {
			logger.Error("Error add new user in db", "repository-AddNewUser()", err)
			return err, false
		}
	}

	logger.Debug("Success add new user in db", "repository-AddNewUser()", nil)
	return nil, false
}

func (r *mainRepository) DeactivateUser(ctx context.Context, chat_id int64) error {
	query := `UPDATE telegram_users
			SET is_active = false
			WHERE chat_id = @chat_id;`
	args := pgx.NamedArgs{
		"chat_id": chat_id,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		logger.Error("Error deactivate user in db", "repository-DeactivateUser()", err)
		return err
	}

	logger.Debug("Success deactivate user in db", "repository-DeactivateUser()", nil)
	return nil
}

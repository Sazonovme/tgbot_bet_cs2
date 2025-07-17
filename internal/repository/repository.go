package repository

import (
	"RushBananaBet/internal/model"
	"RushBananaBet/pkg/logger"
	"context"

	"github.com/jackc/pgx/v5"
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

func (r *mainRepository) CreateMatch(ctx context.Context, match *model.Match) error {

	query := `WITH t AS (
    	SELECT id FROM tournaments WHERE is_active = true LIMIT 1
	)
	INSERT INTO matches (tournament_id, name, team_1, team_2, date)
	SELECT t.id, @name, @team1, @team2, @date
	FROM t`
	args := pgx.NamedArgs{
		"name":  match.Name,
		"team1": match.Team1,
		"team2": match.Team2,
		"date":  match.Date,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		logger.Error("Create match in db error", "repository-createEvent()", err)
		return err
	}
	logger.Debug("Success create match in db", "repository-createEvent()", nil)
	return nil
}

func (r *mainRepository) AddMatchResult(ctx context.Context, result string, match_id int) error {
	query := `UPDATE matches
				SET result = @result
				WHERE id = @match_id;`
	args := pgx.NamedArgs{
		"result":   result,
		"match_id": match_id,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		logger.Error("Add result to event in db error", "repository-addResultToEvent()", err)
		return err
	}
	logger.Debug("Success add result in db", "repository-addResultToEvent()", nil)
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
				AND t.is_active = true`

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

// USER

func (r *mainRepository) GetUserPredictions(ctx context.Context, username string) (*[]model.UserPrediction, error) {

	query := `SELECT
				m.name AS match_name,
				p.prediction			
			FROM predictions p
			JOINT matches m ON p.match_id = m.id
			WHERE
				p.username = @username`
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
		err := sqlRows.Scan(&prediction.Match_Name, &prediction.Prediction)
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
				SELECT @username, @match_id, @prediction
				WHERE NOT EXISTS (
    				SELECT 1 FROM predictions WHERE username = @username AND match_id = @match_id
				);`
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

func (r *mainRepository) AddNewUser(ctx context.Context, user *model.User) error {
	query := `INSERT INTO telegram_users (chat_id, username, first_name, last_name, is_active)
				VALUES (@chat_id, @username, @first_name, @last_name, true)
				ON CONFLICT (chat_id) DO UPDATE SET is_active = true`
	args := pgx.NamedArgs{
		"chat_id":    user.Chat_id,
		"username":   user.Username,
		"first_name": user.First_name,
		"last_name":  user.Last_name,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		logger.Error("Error add new user in db", "repository-AddNewUser()", err)
		return err
	}

	logger.Debug("Success add new user in db", "repository-AddNewUser()", nil)
	return nil
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

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

func (r *mainRepository) CreateEvent(ctx context.Context, event model.Event) error {
	query := `INSERT INTO current_event(name, team_1, team_2, date) VALUES (@name, @team1, @team2, @date)`
	args := pgx.NamedArgs{
		"name":  event.Name,
		"team1": event.Team1,
		"team2": event.Team2,
		"date":  event.Date,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		logger.Error("Create event in db error", "repository-createEvent()", err)
		return err
	}
	logger.Debug("Success create event in db", "repository-createEvent()", nil)
	return nil
}

func (r *mainRepository) AddResultToEvent(ctx context.Context, result string) error {
	query := `INSERT INTO current_event(result) VALUES (@result)`
	args := pgx.NamedArgs{
		"result": result,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		logger.Error("Add result to event in db error", "repository-addResultToEvent()", err)
		return err
	}
	logger.Debug("Success add result in db", "repository-addResultToEvent()", nil)
	return nil
}

func (r *mainRepository) GetEventFinishTable(ctx context.Context) ([]model.FinishTable, error) {
	query := `SELECT 
				cp.username,
				ce.name,
				cp.prediction,
				ce.result
			FROM
				current_event as ce
			INNER JOIN current_predictions as cp 
			ON ce.id = cp.id_event`

	sqlRows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.Error("Get finish table in db error", "repository-getEventFinishTable()", err)
		return nil, err
	}
	defer sqlRows.Close()

	logger.Debug("Success get sql rows for finish table in db", "repository-getEventFinishTable()", nil)

	finishTables := []model.FinishTable{}
	for sqlRows.Next() {
		finishTable := model.FinishTable{}
		err := sqlRows.Scan(&finishTable.Username, &finishTable.Name, &finishTable.Prediction, &finishTable.Result)
		if err != nil {
			logger.Error("Scan finish table in db error", "repository-getEventFinishTable()", err)
			return nil, err
		}
		finishTables = append(finishTables, finishTable)
	}

	logger.Debug("Success get finish table in db", "repository-getEventFinishTable()", nil)
	return finishTables, nil
}

// USER

func (r *mainRepository) GetUserPredictions(ctx context.Context, username string) ([]model.UserPrediction, error) {
	query := `SELECT current_event.name, current_predictions.prediction
			FROM current_predictions AS current_predictions
				LEFT JOIN current_event ON
					current_predictions.id_event = current_event.id
			WHERE
				current_predictions.username = $1`
	sqlRows, err := r.db.Query(ctx, query, username)
	if err != nil {
		logger.Error("Get user predictions in db error", "repository-getUserPredictions()", err)
		return nil, err
	}
	defer sqlRows.Close()

	logger.Debug("Success get sql rows for user predictions in db", "repository-getUserPredictions()", nil)

	predictions := []model.UserPrediction{}
	for sqlRows.Next() {
		prediction := model.UserPrediction{}
		err := sqlRows.Scan(&prediction.NameEvent, &prediction.Prediction)
		if err != nil {
			logger.Error("Scan prediction in db error", "repository-getUserPredictions()", err)
			return nil, err
		}
		predictions = append(predictions, prediction)
	}

	logger.Debug("Success get user predictions in db", "repository-getUserPredictions()", nil)
	return predictions, nil
}

func (r *mainRepository) AddUserPrediction(ctx context.Context, prediction *model.UserPrediction) error {
	query := `INSERT INTO current_predictions(username, id_event, prediction) VALUES (@username, @id_event, @prediction)`
	args := pgx.NamedArgs{
		"username":   prediction.UserName,
		"id_event":   prediction.Id_event,
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

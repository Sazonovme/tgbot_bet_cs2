package service

import (
	"RushBananaBet/internal/model"
	"context"
)

type Service struct {
	Repository Repository
}

type Repository interface {
	CreateEvent(ctx context.Context, event model.Event) error
	AddResultToEvent(ctx context.Context, result string) error
	GetEventFinishTable(ctx context.Context) ([]model.FinishTable, error)
	GetUserPredictions(ctx context.Context, username string) ([]model.UserPrediction, error)
	AddUserPrediction(ctx context.Context, prediction *model.UserPrediction) error
}

func NewService(repo Repository) *Service {
	return &Service{
		Repository: repo,
	}
}

// ADMIN

func (s *Service) CreateEvent(ctx context.Context, event model.Event) error {
	return s.Repository.CreateEvent(ctx, event)
}

func (s *Service) AddResultToEvent(ctx context.Context, result string) error {
	return s.Repository.AddResultToEvent(ctx, result)
}

func (s *Service) GetEventFinishTable(ctx context.Context) ([]model.FinishTable, error) {
	return s.Repository.GetEventFinishTable(ctx)
}

// USER

func (s *Service) GetUserPredictions(ctx context.Context, username string) ([]model.UserPrediction, error) {
	return s.Repository.GetUserPredictions(ctx, username)
}

func (s *Service) AddUserPrediction(ctx context.Context, prediction *model.UserPrediction) error {
	return s.Repository.AddUserPrediction(ctx, prediction)
}

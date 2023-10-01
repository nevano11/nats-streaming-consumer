package service

import (
	"nats-streaming-consumer/internal/entity"
	"nats-streaming-consumer/internal/repository"
)

type DbService struct {
	repository *repository.Repository
}

func NewDbServiceDb(rep *repository.Repository) *DbService {
	return &DbService{
		repository: rep,
	}
}

func (s *DbService) AddNewModel(model entity.Model) (int, error) {
	return s.repository.AddNewModel(model)
}

func (s *DbService) SelectAllModels() ([]entity.Model, error) {
	return s.repository.SelectAllModels()
}

func (s *DbService) SelectModelByUid(uid string) (entity.Model, error) {
	return s.repository.SelectModelByUid(uid)
}

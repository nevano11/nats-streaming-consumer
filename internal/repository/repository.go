package repository

import "nats-streaming-consumer/internal/entity"

type Repository struct {
	ModelCrudRepository
}

func NewRepository(rep ModelCrudRepository) *Repository {
	return &Repository{
		rep,
	}
}

type ModelCrudRepository interface {
	AddNewModel(model entity.Model) (int, error)
	SelectAllModels() ([]entity.Model, error)
	SelectModelByUid(uid string) (entity.Model, error)
}

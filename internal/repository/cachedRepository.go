package repository

import (
	"github.com/sirupsen/logrus"
	"nats-streaming-consumer/internal/entity"
)

type CachedRepository struct {
	baseRepository ModelCrudRepository
	cache          map[string]entity.Model
}

func NewCachedRepository(baseRepository ModelCrudRepository) *CachedRepository {
	return &CachedRepository{
		baseRepository: baseRepository,
		cache:          make(map[string]entity.Model),
	}
}

func (r *CachedRepository) FillCacheFromRepository() error {
	models, err := r.baseRepository.SelectAllModels()
	if err != nil {
		return err
	}
	for _, v := range models {
		r.cache[v.OrderUid] = v
	}
	return nil
}

func (r *CachedRepository) AddNewModel(model entity.Model) (int, error) {
	modelId, err := r.baseRepository.AddNewModel(model)
	if err == nil {
		r.cache[model.OrderUid] = model
	}
	logrus.Debug("Add model on cache")
	return modelId, err
}

func (r *CachedRepository) SelectAllModels() ([]entity.Model, error) {
	return r.baseRepository.SelectAllModels()
}

func (r *CachedRepository) SelectModelByUid(uid string) (entity.Model, error) {
	model, onCache := r.cache[uid]
	if !onCache {
		logrus.Debug("Get model from db")
		byUid, err := r.baseRepository.SelectModelByUid(uid)
		if err != nil {
			return entity.Model{}, err
		}
		r.cache[uid] = byUid
		return byUid, err
	}
	logrus.Debug("Get model from cache")
	return model, nil
}

package repositories

import "flashcard_service/internal/model"

type CategoryRepository interface {
	Insert(userId string, name string) (int64, error)
	UpdateById(userId string, id string, name string) error
	DeleteById(userId string, id string) error
	FindAll(userId string) ([]model.Category, error)
}

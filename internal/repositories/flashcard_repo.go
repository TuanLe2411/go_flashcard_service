package repositories

import (
	"flashcard_service/internal/model"
	"flashcard_service/pkg/objects"
)

type FlashcardRepository interface {
	InsertManyByUserId(userId string, categoryId string, flashcards []objects.CreateFlashcard) error
	FindOneById(userId string, id string) (model.Flashcard, error)
	FindByCategoryId(userId string, categoryId string) ([]model.Flashcard, error)
	DeleteById(userId string, id string) error
	UpdateById(userId string, id string, flashcard model.Flashcard) error
}

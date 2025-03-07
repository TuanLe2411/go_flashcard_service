package repositories

import "flashcard_service/internal/model"

type FlashcardRepository interface {
	InsertManyByUserId(userId string, flashcards []model.Flashcard) error
	FindById(userId string, id string) (model.Flashcard, error)
	DeleteById(userId string, id string) error
	UpdateById(userId string, id string, flashcard model.Flashcard) error
}

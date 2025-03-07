package objects

import "flashcard_service/internal/model"

type UpdateFlashcard struct {
	Name       string
	Content    string
	CategoryId int
}

func (u UpdateFlashcard) ToFlashcard() model.Flashcard {
	return model.Flashcard{
		Name:       u.Name,
		Content:    u.Content,
		CategoryId: u.CategoryId,
	}
}

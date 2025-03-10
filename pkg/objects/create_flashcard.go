package objects

import "flashcard_service/internal/model"

type CreateFlashcard struct {
	Name    string
	Content string
}

func (c CreateFlashcard) ToFlashcard() model.Flashcard {
	return model.Flashcard{
		Name:    c.Name,
		Content: c.Content,
	}
}

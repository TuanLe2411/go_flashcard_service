package objects

import "flashcard_service/internal/model"

type CreateCategory struct {
	Name string
}

func (c CreateCategory) ToCategory() model.Category {
	return model.Category{Name: c.Name}
}

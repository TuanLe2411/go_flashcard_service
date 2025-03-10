package model

import "time"

type Flashcard struct {
	ID         int64      `json:"id,omitempty"`
	Name       string     `json:"name,omitempty"`
	Content    string     `json:"content,omitempty"`
	CategoryId int        `json:"categoryId,omitempty"`
	UserId     int        `json:"userId,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

func (f Flashcard) IsExisted() bool {
	return f.CategoryId > 0
}

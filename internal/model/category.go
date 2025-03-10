package model

import "time"

type Category struct {
	Id        int64      `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	UserID    int        `json:"userId,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

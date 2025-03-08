package repositories_impl

import (
	"flashcard_service/internal/model"
	"flashcard_service/internal/repositories"
	"flashcard_service/pkg/database"
	"time"
)

type CategoryRepositoryImpl struct {
	db database.Database
}

func NewCategoryRepositoryImpl(db database.Database) repositories.CategoryRepository {
	return &CategoryRepositoryImpl{
		db: db,
	}
}

func (c *CategoryRepositoryImpl) Insert(userId string, name string) error {
	_, err := c.db.Exec(
		"INSERT INTO flash_category (name, user_id, created_at) VALUES (?, ?, ?)",
		name,
		userId,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	return err
}

func (c *CategoryRepositoryImpl) FindAll(userId string) ([]model.Category, error) {
	rows, err := c.db.QueryRows(
		"SELECT id, name, user_id, created_at, updated_at FROM flash_category WHERE user_id = ?",
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var category model.Category
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.UserID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *CategoryRepositoryImpl) DeleteById(userId string, id string) error {
	_, err := c.db.Exec(
		"delete from flash_category where user_id = ? and id = ?",
		userId,
		id,
	)
	return err
}

func (c *CategoryRepositoryImpl) UpdateById(userId string, id string, name string) error {
	_, err := c.db.Exec(
		"update flash_category set name = ? where user_id = ? and id = ?",
		name,
		userId,
		id,
	)
	return err
}

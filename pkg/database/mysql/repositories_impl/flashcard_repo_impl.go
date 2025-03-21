package repositories_impl

import (
	"flashcard_service/internal/model"
	"flashcard_service/internal/repositories"
	"flashcard_service/pkg/database"
	"flashcard_service/pkg/objects"
	"fmt"
	"strings"
	"time"
)

type FlashcardRepositoryImpl struct {
	db database.Database
}

func NewFlashcardRepositoryImpl(
	db database.Database,
) repositories.FlashcardRepository {
	return &FlashcardRepositoryImpl{
		db: db,
	}
}

func (f *FlashcardRepositoryImpl) InsertManyByUserId(userId string, categoryId string, flashcards []objects.CreateFlashcard) error {
	valueStrings := make([]string, 0, len(flashcards))
	valueArgs := make([]any, 0, len(flashcards)*3)

	for _, fc := range flashcards {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, fc.Name, fc.Content, categoryId, time.Now().Format("2006-01-02 15:04:05"), userId)
	}

	query := fmt.Sprintf("INSERT INTO flashcard (name, content, category_id, created_at, user_id) VALUES %s",
		strings.Join(valueStrings, ","))

	_, cancel, err := f.db.Exec(query, valueArgs...)
	defer cancel()
	return err
}

func (f *FlashcardRepositoryImpl) FindOneById(userId string, id string) (model.Flashcard, error) {
	sql := "SELECT id, name, content, category_id, created_at, updated_at, user_id FROM flashcard WHERE id = ? and user_id = ?"
	row, cancel, err := f.db.QueryRow(sql, id, userId)
	defer cancel()
	if err != nil {
		return model.Flashcard{}, err
	}
	var flashcard model.Flashcard
	err = row.Scan(
		&flashcard.ID,
		&flashcard.Name,
		&flashcard.Content,
		&flashcard.CategoryId,
		&flashcard.CreatedAt,
		&flashcard.UpdatedAt,
		&flashcard.UserId,
	)
	if err != nil {
		return model.Flashcard{}, err
	}
	return flashcard, nil
}

func (f *FlashcardRepositoryImpl) DeleteById(userId string, id string) error {
	sql := "DELETE FROM flashcard WHERE id = ? and user_id = ?"
	_, cancel, err := f.db.Exec(sql, id, userId)
	defer cancel()
	return err
}

func (f *FlashcardRepositoryImpl) UpdateById(userId string, id string, flashcard model.Flashcard) error {
	sql := "UPDATE flashcard SET name = ?, content = ?, category_id = ?, updated_at = ? WHERE id = ? and user_id = ?"
	_, cancel, err := f.db.Exec(
		sql,
		flashcard.Name,
		flashcard.Content,
		flashcard.CategoryId,
		time.Now().Format("2006-01-02 15:04:05"),
		id,
		userId,
	)
	defer cancel()
	return err
}

func (f *FlashcardRepositoryImpl) FindByCategoryId(userId string, categoryId string) ([]model.Flashcard, error) {
	sql := "SELECT id, name, content, category_id, created_at, updated_at, user_id FROM flashcard WHERE category_id = ? and user_id = ?"
	rows, cancel, err := f.db.QueryRows(sql, categoryId, userId)
	defer cancel()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flashcards []model.Flashcard
	for rows.Next() {
		var flashcard model.Flashcard
		err = rows.Scan(
			&flashcard.ID,
			&flashcard.Name,
			&flashcard.Content,
			&flashcard.CategoryId,
			&flashcard.CreatedAt,
			&flashcard.UpdatedAt,
			&flashcard.UserId,
		)
		if err != nil {
			return nil, err
		}
		flashcards = append(flashcards, flashcard)
	}
	return flashcards, nil
}

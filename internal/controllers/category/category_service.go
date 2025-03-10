package category

import (
	"encoding/json"
	"flashcard_service/internal/model"
	"flashcard_service/pkg/database/redis"
	"strconv"
)

type CategoryService struct {
	r *redis.RedisDatabase
}

func (c *CategoryService) SaveCategoriesToRedisHash(userId string, categories []model.Category) error {
	key := redis.GetCategoriesKey(userId)
	fields := make(map[string]any)
	for _, category := range categories {
		bytes, err := json.Marshal(category)
		if err != nil {
			return err
		}
		strId := strconv.FormatInt(category.Id, 10)
		fields[strId] = string(bytes)
	}
	return c.r.HMSetWithExpiry(key, fields, 300)
}

func (c *CategoryService) SaveCategoryToRedisHash(userId string, categoryId string, category model.Category) error {
	key := redis.GetCategoriesKey(userId)
	bytes, _ := json.Marshal(category)
	return c.r.HSetWithExpiry(key, categoryId, string(bytes), 300)
}

func (c *CategoryService) GetCategoryFromRedisHash(userId string, categoryId string) (model.Category, error) {
	key := redis.GetCategoriesKey(userId)
	value, err := c.r.HGet(key, categoryId)
	if err != nil {
		return model.Category{}, err
	}
	var category model.Category
	err = json.Unmarshal([]byte(value), &category)
	if err != nil {
		return model.Category{}, err
	}
	return category, nil
}

func (c *CategoryService) GetCategoriesFromRedisHash(userId string) ([]model.Category, error) {
	key := redis.GetCategoriesKey(userId)
	values, err := c.r.HGetAll(key)
	if err != nil {
		return nil, err
	}
	categories := make([]model.Category, 0)
	for _, value := range values {
		var category model.Category
		err = json.Unmarshal([]byte(value), &category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *CategoryService) SaveFlashcardsToRedisHash(userId string, categoryId string, flashcards []model.Flashcard) error {
	key := redis.GetFlashcardsKey(userId, categoryId)
	fields := make(map[string]any)
	for _, flashcard := range flashcards {
		bytes, err := json.Marshal(flashcard)
		if err != nil {
			return err
		}
		strId := strconv.FormatInt(flashcard.ID, 10)
		fields[strId] = string(bytes)
	}
	return c.r.HMSetWithExpiry(key, fields, 300)
}

func (c *CategoryService) SaveFlashcardToRedisHash(userId string, categoryId string, flashcard model.Flashcard) error {
	key := redis.GetFlashcardsKey(userId, categoryId)
	bytes, _ := json.Marshal(flashcard)
	return c.r.HSetWithExpiry(key, strconv.FormatInt(flashcard.ID, 10), string(bytes), 300)
}

func (c *CategoryService) GetFlashcardsFromRedisHash(userId string, categoryId string) ([]model.Flashcard, error) {
	key := redis.GetFlashcardsKey(userId, categoryId)
	values, err := c.r.HGetAll(key)
	if err != nil {
		return nil, err
	}
	flashcards := make([]model.Flashcard, 0)
	for _, value := range values {
		var flashcard model.Flashcard
		err = json.Unmarshal([]byte(value), &flashcard)
		if err != nil {
			return nil, err
		}
		flashcards = append(flashcards, flashcard)
	}
	return flashcards, nil
}

func (c *CategoryService) DeleteFlashcardFromRedisHash(userId string, categoryId string, flashcardId string) error {
	key := redis.GetFlashcardsKey(userId, categoryId)
	return c.r.HDel(key, flashcardId)
}

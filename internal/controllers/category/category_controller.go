package category

import (
	"encoding/json"
	"errors"
	"flashcard_service/internal/repositories"
	"flashcard_service/pkg/constant"
	"flashcard_service/pkg/database"
	"flashcard_service/pkg/database/mysql/repositories_impl"
	"flashcard_service/pkg/database/redis"
	"flashcard_service/pkg/objects"
	"flashcard_service/pkg/utils"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
)

type CategoryController struct {
	categoryRepo  repositories.CategoryRepository
	flashcardRepo repositories.FlashcardRepository
	redis         *redis.RedisDatabase
	*CategoryService
}

func NewCategoryController(db database.Database, redis *redis.RedisDatabase) *CategoryController {
	return &CategoryController{
		categoryRepo:  repositories_impl.NewCategoryRepositoryImpl(db),
		flashcardRepo: repositories_impl.NewFlashcardRepositoryImpl(db),
		redis:         redis,
		CategoryService: &CategoryService{
			r: redis,
		},
	}
}

func (c *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	userId := r.Header.Get(constant.UserIdHeader)
	if c.isUserIdInvalid(userId, r, trackingId) {
		return
	}

	var createCategoryRequest objects.CreateCategory
	err := json.NewDecoder(r.Body).Decode(&createCategoryRequest)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "error when parse create category request: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, err)
		return
	}

	id, err := c.categoryRepo.Insert(userId, createCategoryRequest.Name)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "error when create category: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	go func() {
		category := createCategoryRequest.ToCategory()
		category.Id = id
		err = c.CategoryService.SaveCategoryToRedisHash(userId, strconv.FormatInt(id, 10), category)
		if err != nil {
			log.Info().Msg("Failed to update Redis cache after database update: " + err.Error())
		}
	}()

	w.WriteHeader(http.StatusCreated)
}

func (c *CategoryController) GetAllCategory(w http.ResponseWriter, r *http.Request) {
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	userId := r.Header.Get(constant.UserIdHeader)
	if c.isUserIdInvalid(userId, r, trackingId) {
		return
	}

	cachedCategories, err := c.CategoryService.GetCategoriesFromRedisHash(userId)
	if err == nil && len(cachedCategories) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cachedCategories)
		return
	}

	categories, err := c.categoryRepo.FindAll(userId)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "error when get all categories: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	go func() {
		err = c.CategoryService.SaveCategoriesToRedisHash(userId, categories)
		if err != nil {
			log.Info().Msg(err.Error())
			return
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (c *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	if c.isUserIdInvalid(userId, r, trackingId) {
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		msg := "error when delete category: id is empty"
		log.Error().
			Str("trackingId", trackingId).
			Str("error", msg).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, errors.New(msg))
		return
	}

	err := c.categoryRepo.DeleteById(userId, id)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "error when delete category: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}
	go func() {
		err = c.redis.HDel(userId, id)
		if err != nil {
			log.Info().Msg("Failed to update Redis cache after database update: " + err.Error())
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func (c *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	if c.isUserIdInvalid(userId, r, trackingId) {
		return
	}

	var updateCategoryRequest objects.UpdateCategory
	err := json.NewDecoder(r.Body).Decode(&updateCategoryRequest)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "error when parse update category request: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		msg := "error when update category: id is empty"
		log.Error().
			Str("trackingId", trackingId).
			Str("error", msg).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, errors.New(msg))
		return
	}

	err = c.categoryRepo.UpdateById(userId, id, updateCategoryRequest.Name)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "error when update category: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	go func() {
		cachedCategory, err := c.CategoryService.GetCategoryFromRedisHash(userId, id)
		if err == nil {
			cachedCategory.Name = updateCategoryRequest.Name
			err = c.CategoryService.SaveCategoryToRedisHash(userId, id, cachedCategory)
			if err != nil {
				log.Info().Msg("Failed to update Redis cache after database update: " + err.Error())
			}
		} else {
			log.Info().Msg("Cache miss when updating category: " + err.Error())
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func (c *CategoryController) GetFlashcardsByCategoryId(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	if c.isUserIdInvalid(userId, r, trackingId) {
		return
	}

	vars := mux.Vars(r)
	categoryId := vars["category_id"]
	if len(categoryId) == 0 {
		msg := "error when get flashcards by category id: category id is empty"
		log.Error().
			Str("trackingId", trackingId).
			Str("error", msg).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, errors.New(msg))
		return
	}

	cachedFlashcards, err := c.CategoryService.GetFlashcardsFromRedisHash(userId, categoryId)
	if err == nil && len(cachedFlashcards) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cachedFlashcards)
		return
	}

	flashcards, err := c.flashcardRepo.FindByCategoryId(userId, categoryId)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "error when get flashcards by category id: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	go func() {
		err = c.CategoryService.SaveFlashcardsToRedisHash(userId, categoryId, flashcards)
		if err != nil {
			log.Info().Msg(err.Error())
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flashcards)
}

func (c *CategoryController) CreateNewFlashcards(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	if c.isUserIdInvalid(userId, r, trackingId) {
		return
	}

	vars := mux.Vars(r)
	categoryId := vars["category_id"]
	if len(categoryId) == 0 {
		msg := "error when create flashcards: category id is empty"
		log.Error().
			Str("trackingId", trackingId).
			Str("error", msg).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, errors.New(msg))
		return
	}

	var createFlashcardsRequest []objects.CreateFlashcard
	err := json.NewDecoder(r.Body).Decode(&createFlashcardsRequest)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "error when parse create flashcards request: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, err)
		return
	}

	err = c.flashcardRepo.InsertManyByUserId(userId, categoryId, createFlashcardsRequest)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "error when create flashcards: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	go func() {
		flashcards, err := c.flashcardRepo.FindByCategoryId(userId, categoryId)
		if err != nil {
			log.Info().Msg("Failed to fetch flashcards for cache update: " + err.Error())
			return
		}
		err = c.CategoryService.SaveFlashcardsToRedisHash(userId, categoryId, flashcards)
		if err != nil {
			log.Info().Msg("Failed to update Redis cache after creating flashcards: " + err.Error())
		}
	}()

	w.WriteHeader(http.StatusCreated)
}

func (c *CategoryController) DeleteFlashcard(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	if c.isUserIdInvalid(userId, r, trackingId) {
		return
	}

	vars := mux.Vars(r)
	categoryId := vars["category_id"]
	flashcardId := vars["flashcard_id"]
	if c.isFlashcardAndCategoryInvalid(categoryId, flashcardId, r, trackingId) {
		return
	}

	err := c.flashcardRepo.DeleteById(userId, flashcardId)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "error when delete flashcard: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	go func() {
		err = c.CategoryService.DeleteFlashcardFromRedisHash(userId, categoryId, flashcardId)
		if err != nil {
			log.Info().Msg("Failed to delete Redis cache after database delete: " + err.Error())
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func (c *CategoryController) UpdateFlashcard(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	if c.isUserIdInvalid(userId, r, trackingId) {
		return
	}

	vars := mux.Vars(r)
	categoryId := vars["category_id"]
	flashcardId := vars["flashcard_id"]
	if c.isFlashcardAndCategoryInvalid(categoryId, flashcardId, r, trackingId) {
		return
	}

	var updateFlashcardRequest objects.UpdateFlashcard
	err := json.NewDecoder(r.Body).Decode(&updateFlashcardRequest)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "error when parse update flashcard request: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, err)
		return
	}

	flashcard := updateFlashcardRequest.ToFlashcard()
	err = c.flashcardRepo.UpdateById(userId, flashcardId, flashcard)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "error when update flashcard: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	go func() {
		flashcardIdInt, err := strconv.ParseInt(flashcardId, 10, 64)
		if err != nil {
			log.Info().Msg("Failed to parse flashcard ID: " + err.Error())
			return
		}
		flashcard.ID = flashcardIdInt
		err = c.CategoryService.SaveFlashcardToRedisHash(userId, categoryId, flashcard)
		if err != nil {
			log.Info().Msg("Failed to update Redis cache after database update: " + err.Error())
		}

	}()

	w.WriteHeader(http.StatusOK)
}

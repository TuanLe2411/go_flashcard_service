package category

import (
	"encoding/json"
	"flashcard_service/internal/repositories"
	"flashcard_service/pkg/constant"
	"flashcard_service/pkg/database"
	"flashcard_service/pkg/database/mysql/repositories_impl"
	"flashcard_service/pkg/database/redis"
	"flashcard_service/pkg/objects"
	"flashcard_service/pkg/utils"
	"log"
	"net/http"
	"strconv"

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
	userId := r.Header.Get(constant.UserIdHeader)
	if len(userId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	var createCategoryRequest objects.CreateCategory
	err := json.NewDecoder(r.Body).Decode(&createCategoryRequest)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	id, err := c.categoryRepo.Insert(userId, createCategoryRequest.Name)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	go func() {
		category := createCategoryRequest.ToCategory()
		category.Id = id
		err = c.CategoryService.SaveCategoryToRedisHash(userId, strconv.FormatInt(id, 10), category)
		if err != nil {
			log.Println("Failed to update Redis cache after database update: ", err)
		}
	}()

	w.WriteHeader(http.StatusCreated)
}

func (c *CategoryController) GetAllCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	if len(userId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
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
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	go func() {
		err = c.CategoryService.SaveCategoriesToRedisHash(userId, categories)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (c *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	if len(userId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	err := c.categoryRepo.DeleteById(userId, id)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}
	go func() {
		err = c.redis.HDel(userId, id)
		if err != nil {
			log.Println("Failed to update Redis cache after database update: ", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func (c *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	if len(userId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	var updateCategoryRequest objects.UpdateCategory
	err := json.NewDecoder(r.Body).Decode(&updateCategoryRequest)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	err = c.categoryRepo.UpdateById(userId, id, updateCategoryRequest.Name)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	go func() {
		cachedCategory, err := c.CategoryService.GetCategoryFromRedisHash(userId, id)
		if err == nil {
			cachedCategory.Name = updateCategoryRequest.Name
			err = c.CategoryService.SaveCategoryToRedisHash(userId, id, cachedCategory)
			if err != nil {
				log.Println("Failed to update Redis cache after database update: ", err)
			}
		} else {
			log.Println("Cache miss when updating category: ", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func (c *CategoryController) GetFlashcardsByCategoryId(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	if len(userId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	vars := mux.Vars(r)
	categoryId := vars["category_id"]
	if len(categoryId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	cachedFlashcards, err := c.CategoryService.GetFlashcardsFromRedisHash(userId, categoryId)
	if err == nil && len(cachedFlashcards) > 0 {
		log.Println("Cache hit")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cachedFlashcards)
		return
	}

	flashcards, err := c.flashcardRepo.FindByCategoryId(userId, categoryId)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	go func() {
		err = c.CategoryService.SaveFlashcardsToRedisHash(userId, categoryId, flashcards)
		if err != nil {
			log.Println(err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flashcards)
}

func (c *CategoryController) CreateNewFlashcards(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	if len(userId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	vars := mux.Vars(r)
	categoryId := vars["category_id"]
	if len(categoryId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	var createFlashcardsRequest []objects.CreateFlashcard
	err := json.NewDecoder(r.Body).Decode(&createFlashcardsRequest)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	err = c.flashcardRepo.InsertManyByUserId(userId, categoryId, createFlashcardsRequest)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *CategoryController) DeleteFlashcard(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	if len(userId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	vars := mux.Vars(r)
	categoryId := vars["category_id"]
	flashcardId := vars["flashcard_id"]
	if len(categoryId) == 0 || len(flashcardId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	err := c.flashcardRepo.DeleteById(userId, flashcardId)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	go func() {
		err = c.CategoryService.DeleteFlashcardFromRedisHash(userId, categoryId, flashcardId)
		if err != nil {
			log.Println("Failed to delete Redis cache after database delete: ", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func (c *CategoryController) UpdateFlashcard(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	if len(userId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	vars := mux.Vars(r)
	categoryId := vars["category_id"]
	flashcardId := vars["flashcard_id"]
	if len(categoryId) == 0 || len(flashcardId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	var updateFlashcardRequest objects.UpdateFlashcard
	err := json.NewDecoder(r.Body).Decode(&updateFlashcardRequest)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	flashcard := updateFlashcardRequest.ToFlashcard()
	err = c.flashcardRepo.UpdateById(userId, flashcardId, flashcard)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	go func() {
		flashcardIdInt, err := strconv.ParseInt(flashcardId, 10, 64)
		if err != nil {
			log.Println("Failed to parse flashcard ID:", err)
			return
		}
		flashcard.ID = flashcardIdInt
		err = c.CategoryService.SaveFlashcardToRedisHash(userId, categoryId, flashcard)
		if err != nil {
			log.Println("Failed to update Redis cache after database update: ", err)
		}

	}()

	w.WriteHeader(http.StatusOK)
}

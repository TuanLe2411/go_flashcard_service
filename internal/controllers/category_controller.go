package controllers

import (
	"encoding/json"
	"flashcard_service/internal/repositories"
	"flashcard_service/pkg/constant"
	"flashcard_service/pkg/database"
	"flashcard_service/pkg/database/mysql/repositories_impl"
	"flashcard_service/pkg/objects"
	"flashcard_service/pkg/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type CategoryController struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryController(db database.Database) *CategoryController {
	return &CategoryController{
		categoryRepo: repositories_impl.NewCategoryRepositoryImpl(db),
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
	err = c.categoryRepo.Insert(userId, createCategoryRequest.Name)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *CategoryController) GetAllCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	if len(userId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	categories, err := c.categoryRepo.FindAll(userId)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

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
	w.WriteHeader(http.StatusOK)
}

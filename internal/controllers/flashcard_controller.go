package controllers

import (
	"encoding/json"
	"flashcard_service/internal/model"
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

type FlashcardController struct {
	flashcardRepo repositories.FlashcardRepository
}

func NewFlashcardController(
	db database.Database,
) *FlashcardController {
	return &FlashcardController{
		flashcardRepo: repositories_impl.NewFlashcardRepositoryImpl(db),
	}
}

func (f *FlashcardController) CreateFlashcards(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	if len(userId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	var flashcards []model.Flashcard
	err := json.NewDecoder(r.Body).Decode(&flashcards)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}
	err = f.flashcardRepo.InsertManyByUserId(userId, flashcards)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (f *FlashcardController) GetByID(w http.ResponseWriter, r *http.Request) {
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
	flashcard, err := f.flashcardRepo.FindById(userId, id)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}
	if !flashcard.IsExisted() {
		utils.SetHttpReponseError(r, utils.ErrNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flashcard)
}

func (f *FlashcardController) UpdateByID(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get(constant.UserIdHeader)
	if len(userId) == 0 {
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	var updateFlashcard objects.UpdateFlashcard
	err := json.NewDecoder(r.Body).Decode(&updateFlashcard)
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

	flashcard := updateFlashcard.ToFlashcard()
	err = f.flashcardRepo.UpdateById(userId, id, flashcard)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (f *FlashcardController) DeleteByID(w http.ResponseWriter, r *http.Request) {
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

	err := f.flashcardRepo.DeleteById(userId, id)
	if err != nil {
		log.Println(err)
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}
}

package drivers

import (
	"flashcard_service/internal/app_log"
	"flashcard_service/internal/controllers/app"
	"flashcard_service/internal/controllers/category"
	"flashcard_service/internal/middleware"
	"flashcard_service/pkg"
	"flashcard_service/pkg/database/mysql"
	"flashcard_service/pkg/database/redis"

	"github.com/rs/zerolog/log"

	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const apiV1Prefix = "/api/v1"

const CategoryControllerPrefix = "/category"
const CreateCategory = ""
const GetCateforyByID = "/{id}"
const GetAllCategories = ""
const UpdateCategoryByID = "/{id}"
const DeleteCategoryByID = "/{id}"
const GetFlashcards = "/{category_id}/flashcards"
const CreateNewFlashcards = "/{category_id}/flashcards"
const DeleteFlashcard = "/{category_id}/flashcards/{flashcard_id}"
const UpdateFlashcard = "/{category_id}/flashcards/{flashcard_id}"

const HeathCheck = "/health"

func Run() {
	pkg.LoadConfig()
	app_log.InitLogger()
	sqlDb := mysql.NewMySql()
	err := sqlDb.Connect()
	if err != nil {
		log.Fatal().Msg("Error when connect to db: " + err.Error())
	}
	log.Info().Msg("Connect to db successfully")

	redis := redis.NewRedisClient()
	err = redis.Connect()
	if err != nil {
		log.Fatal().Msg("Error when connect to redis: " + err.Error())
		panic(err)
	}
	log.Info().Msg("Connect to redis successfully")

	router := mux.NewRouter()
	appController := &app.AppController{}
	router.HandleFunc(HeathCheck, appController.HeathCheck).Methods(http.MethodGet)

	router.Use(
		middleware.XssProtectionMiddleware,
		middleware.CorsMiddleware,
		middleware.MonitorMiddleware,
		middleware.ErrorHandlerMiddleware,
	)

	baseRouter := router.PathPrefix(apiV1Prefix).Subrouter()

	categoryController := category.NewCategoryController(sqlDb, redis)
	categoryRouter := baseRouter.PathPrefix(CategoryControllerPrefix).Subrouter()
	categoryRouter.HandleFunc(CreateCategory, categoryController.CreateCategory).Methods(http.MethodPost)
	categoryRouter.HandleFunc(GetAllCategories, categoryController.GetAllCategory).Methods(http.MethodGet)
	categoryRouter.HandleFunc(DeleteCategoryByID, categoryController.DeleteCategory).Methods(http.MethodDelete)
	categoryRouter.HandleFunc(UpdateCategoryByID, categoryController.UpdateCategory).Methods(http.MethodPut)
	categoryRouter.HandleFunc(CreateNewFlashcards, categoryController.CreateNewFlashcards).Methods(http.MethodPost)
	categoryRouter.HandleFunc(GetFlashcards, categoryController.GetFlashcardsByCategoryId).Methods(http.MethodGet)
	categoryRouter.HandleFunc(DeleteFlashcard, categoryController.DeleteFlashcard).Methods(http.MethodDelete)
	categoryRouter.HandleFunc(UpdateFlashcard, categoryController.UpdateFlashcard).Methods(http.MethodPut)

	log.Info().Msg("Server is running on port: " + os.Getenv("SERVER_PORT"))
	http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), router)
}

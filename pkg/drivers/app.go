package drivers

import (
	"flashcard_service/internal/controllers"
	"flashcard_service/internal/middleware"
	"flashcard_service/pkg"
	"flashcard_service/pkg/database/mysql"
	"flashcard_service/pkg/database/redis"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const apiV1Prefix = "/api/v1"

const FlashcardControllerPrefix = "/flashcard"
const CreateFlashcards = ""
const GetFlashcardByID = "/{id}"
const UpdateFlashcardByID = "/{id}"
const DeleteFlashcardByID = "/{id}"

const CategoryControllerPrefix = "/category"
const CreateCategory = ""
const GetCateforyByID = "/{id}"
const GetAllCategories = ""
const UpdateCategoryByID = "/{id}"
const DeleteCategoryByID = "/{id}"

func Run() {
	pkg.LoadConfig()
	sqlDb := mysql.NewMySql()
	err := sqlDb.Connect()
	if err != nil {
		log.Println("Error when connect to db: " + err.Error())
		panic(err)
	}
	log.Println("Connect to db successfully")

	redis := redis.NewRedisClient()
	err = redis.Connect()
	if err != nil {
		log.Println("Error when connect to redis: " + err.Error())
		panic(err)
	}
	log.Println("Connect to redis successfully")

	router := mux.NewRouter()

	router.Use(
		middleware.XssProtectionMiddleware,
		middleware.CorsMiddleware,
		middleware.MonitorMiddleware,
		middleware.ErrorHandlerMiddleware,
	)

	baseRouter := router.PathPrefix(apiV1Prefix).Subrouter()

	flashcardController := controllers.NewFlashcardController(sqlDb)
	flashcardRouter := baseRouter.PathPrefix(FlashcardControllerPrefix).Subrouter()
	flashcardRouter.HandleFunc(CreateFlashcards, flashcardController.CreateFlashcards).Methods(http.MethodPost)
	flashcardRouter.HandleFunc(GetFlashcardByID, flashcardController.GetByID).Methods(http.MethodGet)
	flashcardRouter.HandleFunc(UpdateFlashcardByID, flashcardController.UpdateByID).Methods(http.MethodPut)
	flashcardRouter.HandleFunc(DeleteFlashcardByID, flashcardController.DeleteByID).Methods(http.MethodDelete)

	categoryController := controllers.NewCategoryController(sqlDb)
	categoryRouter := baseRouter.PathPrefix(CategoryControllerPrefix).Subrouter()
	categoryRouter.HandleFunc(CreateCategory, categoryController.CreateCategory).Methods(http.MethodPost)
	categoryRouter.HandleFunc(GetAllCategories, categoryController.GetAllCategory).Methods(http.MethodGet)
	categoryRouter.HandleFunc(DeleteCategoryByID, categoryController.DeleteCategory).Methods(http.MethodDelete)
	categoryRouter.HandleFunc(UpdateCategoryByID, categoryController.UpdateCategory).Methods(http.MethodPut)

	log.Println("Server is running on port: " + os.Getenv("SERVER_PORT"))
	http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), router)
}

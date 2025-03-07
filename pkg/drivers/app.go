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

const flashcardControllerPrefix = "/flashcard"
const creates = ""
const getByID = "/{id}"
const UpdateByID = "/{id}"
const deleteByID = "/{id}"

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

	flashcardRouter := baseRouter.PathPrefix(flashcardControllerPrefix).Subrouter()
	flashcardRouter.HandleFunc(creates, flashcardController.CreateFlashcards).Methods(http.MethodPost)
	flashcardRouter.HandleFunc(getByID, flashcardController.GetByID).Methods(http.MethodGet)
	flashcardRouter.HandleFunc(UpdateByID, flashcardController.UpdateByID).Methods(http.MethodPut)
	flashcardRouter.HandleFunc(deleteByID, flashcardController.DeleteByID).Methods(http.MethodDelete)

	log.Println("Server is running on port: " + os.Getenv("SERVER_PORT"))
	http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), router)
}

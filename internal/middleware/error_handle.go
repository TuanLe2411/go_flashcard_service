package middleware

import (
	"encoding/json"
	"flashcard_service/pkg/constant"
	"flashcard_service/pkg/objects"
	"flashcard_service/pkg/utils"
	"net/http"
)

func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		if err, ok := r.Context().Value(constant.AppErrorContextKey).(utils.AppError); ok {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(objects.ApiResponse{
				Code:    err.Code,
				Message: err.Message,
			})
		}
	})
}

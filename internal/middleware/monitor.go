package middleware

import (
	"context"
	"flashcard_service/internal/app_log"
	"flashcard_service/pkg/constant"
	"flashcard_service/pkg/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func MonitorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trackingId := uuid.New().String()
		ctx := context.WithValue(r.Context(), constant.TrackingIdContextKey, trackingId)
		r = r.WithContext(ctx)
		start := time.Now()

		next.ServeHTTP(w, r)

		isError := false
		responseCode := http.StatusOK
		errorMessage := ""
		endTime := time.Now()
		durationMs := endTime.Sub(start).Milliseconds()
		if err, ok := r.Context().Value(constant.AppErrorContextKey).(utils.AppError); ok {
			isError = true
			responseCode = err.Code
			errorMessage = err.ErrorMessage
		}
		app_log.LogKQI(app_log.KQI{
			TrackingId:   trackingId,
			LogType:      app_log.API,
			HttpMethod:   r.Method,
			HttpPath:     r.URL.Path,
			FunctionName: "",
			OriginalName: "",
			ServiceName:  "User Service",
			Description:  "",
			StartTme:     start,
			EndTime:      endTime,
			DurationMs:   durationMs,
			IsError:      isError,
			ErrorMessage: errorMessage,
			ResponseCode: responseCode,
		})
	})
}

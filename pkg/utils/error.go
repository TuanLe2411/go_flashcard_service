package utils

import (
	"context"
	"flashcard_service/pkg/constant"
	"net/http"
)

type AppError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e AppError) Error() string {
	return e.Message
}

// Một số lỗi cụ thể
var (
	ErrNotFound     = AppError{Code: http.StatusNotFound, Message: "Không tìm thấy tài nguyên"}
	ErrBadRequest   = AppError{Code: http.StatusBadRequest, Message: "Yêu cầu không hợp lệ"}
	ErrServerError  = AppError{Code: http.StatusInternalServerError, Message: "Lỗi máy chủ"}
	ErrUnAuthorized = AppError{Code: http.StatusUnauthorized, Message: "Không có quyền truy cập"}
)

func SetHttpReponseError(r *http.Request, err AppError) {
	ctx := context.WithValue(r.Context(), constant.AppErrorContextKey, err)
	*r = *r.WithContext(ctx)
}

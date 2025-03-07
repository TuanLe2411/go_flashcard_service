package utils

import (
	"flashcard_service/pkg/constant"
	"net/http"
)

func ChainMiddlewares(handler http.Handler, middlewares ...constant.Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i].Do(handler)
	}
	return handler
}

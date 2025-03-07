package constant

import "net/http"

type Middleware interface {
	Do(http.Handler) http.Handler
}

type contextKey string

type UserAction string

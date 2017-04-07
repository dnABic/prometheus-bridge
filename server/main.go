package server

import (
	"context"
	"net/http"
)

type ServerHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request)

func HandleWithContext(ctx context.Context, f ServerHandler) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f(ctx, w, r)
	})
}

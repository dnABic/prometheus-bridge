package server

import (
	"context"
	"net/http"
)

type ContextHandler func(context.Context, http.ResponseWriter, *http.Request)

func HandleWithContext(ctx context.Context, f ContextHandler) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f(ctx, w, r)
	})
}

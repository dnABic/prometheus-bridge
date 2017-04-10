package server

import (
	"prometheus-bridge/messaging"

	"golang.org/x/net/context"
)

type appContext struct {
	context.Context

	stream messaging.Stream
}

type key int

const streamKey key = 0

func NewContext(parent context.Context, stream messaging.Stream) context.Context {
	return &appContext{parent, stream}
}

func (ctx *appContext) Value(key interface{}) interface{} {
	if key == streamKey {
		return ctx.stream
	}

	return ctx.Context.Value(key)
}

func MessagingStream(ctx context.Context) (messaging.Stream, bool) {
	str, ok := ctx.Value(streamKey).(messaging.Stream)

	return str, ok
}

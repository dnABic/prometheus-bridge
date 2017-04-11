package server

import (
	"context"
	"prometheus-bridge/messaging"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockStream struct {
	messaging.Stream
}

func TestCreateNewContext(t *testing.T) {
	var stream messaging.Stream = (*mockStream)(nil)

	ctx := NewContext(context.Background(), stream)
	s, ok := MessagingStream(ctx)

	assert.True(t, ok)
	assert.Equal(t, stream, s)
}

func TestGetStreamFromContectHeirarchy(t *testing.T) {
	var stream messaging.Stream = (*mockStream)(nil)

	parent := NewContext(context.Background(), stream)
	ctx := context.WithValue(parent, "Test", "test")
	s, ok := MessagingStream(ctx)

	assert.True(t, ok)
	assert.Equal(t, stream, s)
}

func TestContextGetValue(t *testing.T) {
	var stream messaging.Stream = (*mockStream)(nil)

	parent := NewContext(context.Background(), stream)
	ctx := context.WithValue(parent, "Test", "test")
	s := ctx.Value("Test")

	assert.Equal(t, "test", s)
}

func TestContextReturnsNilForNonExistingKeys(t *testing.T) {
	var stream messaging.Stream = (*mockStream)(nil)

	ctx := NewContext(context.Background(), stream)
	s := ctx.Value("Test")

	assert.Nil(t, s)
}

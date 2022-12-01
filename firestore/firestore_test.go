package firestore

import (
	"context"
	"testing"
)

var ctx = context.Background()

func TestFirestoreCache(t *testing.T) {
	cache, close := New[Sample](ctx, &FirestoreConfig{
		ProjectID:  "holy-diver-297719",
		Collection: "ezcx-cache-testing",
		Document:   "maps",
	})
	defer close()
	cache.Listen(ctx)
	t.Log(cache.Get("sample-1"))
	t.Log(cache.Get("sample-2"))
}

type Sample struct {
	String1 string `json:"string-1,omitempty" firestore:"string-1,omitempty"`
	String2 string `json:"string-2,omitempty" firestore:"string-2,omitempty"`
	String3 string `json:"string-3,omitempty" firestore:"string-3,omitempty"`
	String4 string `json:"string-4,omitempty" firestore:"string-4,omitempty"`
}

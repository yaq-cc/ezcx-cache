package firestore

import (
	"context"
	"errors"
	"log"
	"sync"

	fs "cloud.google.com/go/firestore"
	cache "github.com/yaq-cc/ezcx-cache"
)

var (
	ErrValueTypeMismatch = errors.New("value type mismatch")
)

// Firestore's snapshot data returns data in a map[string]any structure
// Safe to assume we can treat integer keys as strings unless mathematical
// operations are necessary.  Handle those separately.  
type FirestoreCache[V any] struct {
	cfg    *FirestoreConfig
	client *fs.Client
	cache  *cache.Cache[string, V]
}

func New[V any](ctx context.Context, cfg *FirestoreConfig) (*FirestoreCache[V], func() error) {
	c := new(FirestoreCache[V])
	c.cfg = cfg
	client, err := fs.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		log.Fatal(err)
	}
	c.client = client
	c.cache = cache.New[string, V]()
	return c, c.client.Close
}

func (c *FirestoreCache[V]) Get(key string) (V, bool) {
	return c.cache.Get(key)
}

func (c *FirestoreCache[V]) Listen(ctx context.Context) {
	changes := c.client.Collection(c.cfg.Collection).Doc(c.cfg.Document).Snapshots(ctx)
	var once sync.Once
	ready := make(chan struct{})
	readyFunc := func() {
		ready <- struct{}{}
		close(ready)
	}

	go func() {
		for {
			change, err := changes.Next()
			if err != nil {
				log.Fatal(err)
			}
			for key, value := range change.Data() {
				asserted, ok := value.(V)
				if !ok {
					log.Fatal(ErrValueTypeMismatch)
				}
				c.cache.Set(key, asserted)
			}
	
			once.Do(readyFunc)
		}
	}()

	<-ready
}

type FirestoreConfig struct {
	ProjectID  string
	Collection string
	Document   string
}

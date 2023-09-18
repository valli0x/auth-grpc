package inmem

import (
	"context"
	"sync"

	"github.com/valli0x/auth-grpc/storage"
)

type InmemStorage struct {
	sync.RWMutex
	db map[string]*storage.Entry
}

func NewInmem() (storage.Storage, error) {
	return &InmemStorage{
		db: make(map[string]*storage.Entry, 100), // just 100, I hope you don't want to add more entries
	}, nil
}

func (i *InmemStorage) Put(ctx context.Context, entry *storage.Entry) error {
	i.Lock()
	defer i.Unlock()

	return i.PutInternal(ctx, entry)
}

func (i *InmemStorage) PutInternal(ctx context.Context, entry *storage.Entry) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	i.db[entry.Key] = entry
	return nil
}

func (i *InmemStorage) Get(ctx context.Context, key string) (*storage.Entry, error) {
	i.RLock()
	defer i.RUnlock()

	return i.GetInternal(ctx, key)
}

func (i *InmemStorage) GetInternal(ctx context.Context, key string) (*storage.Entry, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if entry, ok := i.db[key]; ok {
		return entry, nil
	}

	return nil, nil
}

func (i *InmemStorage) Delete(ctx context.Context, key string) error {
	i.Lock()
	defer i.Unlock()

	return i.DeleteInternal(ctx, key)
}

func (i *InmemStorage) DeleteInternal(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	delete(i.db, key)
	return nil
}

func (i *InmemStorage) List(ctx context.Context) ([]string, error) {
	i.RLock()
	defer i.RUnlock()

	return i.ListInternal(ctx)
}

func (i *InmemStorage) ListInternal(ctx context.Context) ([]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	keys := make([]string, 0, 100)

	for key, _ := range i.db {
		keys = append(keys, key)
	}

	return keys, nil
}
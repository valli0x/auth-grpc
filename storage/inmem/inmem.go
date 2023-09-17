package inmem

import (
	"context"
	"sync"

	"github.com/valli0x/auth-grpc/storage"
)

type InmemStorage struct {
	sync.RWMutex
	db map[string]interface{}
}

func NewInmem() (storage.Storage, error) {
	return &InmemStorage{
		db: make(map[string]interface{}, 100), // just 100, I dont want thinking about it
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

	return nil, nil
}

func (i *InmemStorage) GetInternal(ctx context.Context, key string) (*storage.Entry, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if entry, ok := i.db[key]; ok {
		return entry.(*storage.Entry), nil
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

func (i *InmemStorage) List(ctx context.Context) ([]*storage.Entry, error) {
	i.RLock()
	defer i.RUnlock()

	return i.ListInternal(ctx)
}

func (i *InmemStorage) ListInternal(ctx context.Context) ([]*storage.Entry, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	listEntry := []*storage.Entry{}

	for _, entry := range i.db {
		listEntry = append(listEntry, entry.(*storage.Entry))
	}

	return listEntry, nil
}
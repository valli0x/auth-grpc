package storage

import "context"

type Storage interface {
	Put(ctx context.Context, entry *Entry) error

	Get(ctx context.Context, key string) (*Entry, error)

	Delete(ctx context.Context, key string) error

	List(ctx context.Context) ([]*Entry, error)
}

type Entry struct {
	Key string
	Val interface{}
}
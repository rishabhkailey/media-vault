package mediasearch

import (
	"context"
)

// todo move get user media here
type Service interface {
	CreateOne(context.Context, Model) (int64, error)
	CreateMany(context.Context, []Model) (int64, error)
	DeleteOne(context.Context, DeleteOneCommand) (int64, error)
	DeleteMany(context.Context, DeleteManyCommand) (int64, error)
	Search(context.Context, MediaSearchQuery) ([]Model, error)
	SearchGetMediaIDsOnly(context.Context, MediaSearchQuery) ([]uint, error)
}

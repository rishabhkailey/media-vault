package mediasearch

import (
	"context"
)

// todo move get user media here
type Service interface {
	CreateOne(context.Context, Model) (int64, error)
	CreateMany(context.Context, []Model) (int64, error)
	Search(context.Context, MediaSearchQuery) ([]Model, error)
	SearchGetMediaIDsOnly(context.Context, MediaSearchQuery) ([]uint, error)
}

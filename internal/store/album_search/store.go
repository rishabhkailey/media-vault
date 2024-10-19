package albumsearch

import "context"

type Store interface {
	SearchGetIDsOrderByCreatedAt(ctx context.Context, query, userID string, page, perPage int64, sort string) ([]uint, error)
	SearchGetIDsOrderByUploadedAt(ctx context.Context, query, userID string, page, perPage int64, sort string) ([]uint, error)
	Insert(ctx context.Context, documents []Index) (int64, error)
	Delete(ctx context.Context, ids []string) (int64, error)
	MonitorTask(taskID int64, message string) error
}

package mediasearchimpl

import (
	"context"
	"fmt"
	"strconv"

	"github.com/meilisearch/meilisearch-go"
	mediasearch "github.com/rishabhkailey/media-service/internal/services/mediaSearch"
)

type Service struct {
	store store
}

var _ mediasearch.Service = (*Service)(nil)

func NewService(ms *meilisearch.Client) (mediasearch.Service, error) {
	store, err := newMeiliSearchStore(ms)
	if err != nil {
		return nil, err
	}
	return &Service{
		store: store,
	}, nil
}

func (s *Service) CreateOne(ctx context.Context, mediaSearchData mediasearch.Model) (taskID int64, err error) {
	taskID, err = s.store.Insert(ctx, []mediasearch.Model{mediaSearchData})
	go s.store.MonitorTask(taskID, fmt.Sprintf("[MediaSearch] add document primaryKey=%v", mediaSearchData.MediaID))
	return
}

func (s *Service) CreateMany(ctx context.Context, mediaSearchDataList []mediasearch.Model) (taskID int64, err error) {
	taskID, err = s.store.Insert(ctx, mediaSearchDataList)
	ids := ""
	for _, mediaSearchData := range mediaSearchDataList {
		ids += fmt.Sprintf(" %d", mediaSearchData.MediaID)
	}
	go s.store.MonitorTask(taskID, fmt.Sprintf("[MediaSearch] add document IDs=%v", ids))
	return
}

func (s *Service) DeleteOne(ctx context.Context, cmd mediasearch.DeleteOneCommand) (taskID int64, err error) {
	taskID, err = s.store.Delete(ctx, []string{strconv.FormatUint(uint64(cmd.MediaID), 10)})
	go s.store.MonitorTask(taskID, fmt.Sprintf("[MediaSearch] delete document primaryKey=%v", cmd.MediaID))
	return
}

func (s *Service) DeleteMany(ctx context.Context, cmd mediasearch.DeleteManyCommand) (taskID int64, err error) {
	var ids []string
	for _, id := range cmd.MediaIDs {
		ids = append(ids, strconv.FormatUint(uint64(id), 10))
	}
	taskID, err = s.store.Delete(ctx, ids)
	go s.store.MonitorTask(taskID, fmt.Sprintf("[MediaSearch] delete documents primaryKey=%v", cmd.MediaIDs))
	return
}

func (s *Service) Search(ctx context.Context, query mediasearch.MediaSearchQuery) ([]mediasearch.Model, error) {
	return s.store.Search(ctx, query)
}

func (s *Service) SearchGetMediaIDsOnly(ctx context.Context, query mediasearch.MediaSearchQuery) ([]uint, error) {
	return s.store.SearchGetMediaIDsOnly(ctx, query)
}

package albumsearchimpl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	albumsearch "github.com/rishabhkailey/media-vault/internal/store/album_search"
	"github.com/sirupsen/logrus"
)

type meiliSearchStore struct {
	cli   meilisearch.ServiceManager
	index meilisearch.IndexManager
}

var _ albumsearch.Store = (*meiliSearchStore)(nil)

func NewMeiliSearchStore(cli meilisearch.ServiceManager) (albumsearch.Store, error) {
	albumIndex := cli.Index(albumsearch.GetIndexName())

	// todo move this to migration scripts
	searchableAttributes := albumsearch.GetSearchableAttributes()
	resp, err := albumIndex.UpdateSearchableAttributes(&searchableAttributes)
	if err != nil {
		return nil, fmt.Errorf("[newMeiliSearchStore] update searchable attributes failed: %w", err)
	}
	task, err := cli.WaitForTask(resp.TaskUID, 0)
	if err == nil && task.Status == meilisearch.TaskStatusFailed {
		err = errors.New(task.Error.Message)
	}
	if err != nil {
		return nil, fmt.Errorf("[newMeiliSearchStore] update searchable attributes failed: %w", err)
	} else {
		logrus.Info("[newMeiliSearchStore] succesfuly update searchable attribute")
	}

	// UpdateFilterableAttributes
	filterableAttributes := albumsearch.GetFilterableAttributes()
	resp, err = albumIndex.UpdateFilterableAttributes(&filterableAttributes)
	if err != nil {
		return nil, fmt.Errorf("[newMeiliSearchStore] update searchable attributes failed: %w", err)
	}
	task, err = cli.WaitForTask(resp.TaskUID, 0)
	if err == nil && task.Status == meilisearch.TaskStatusFailed {
		err = errors.New(task.Error.Message)
	}
	if err != nil {
		return nil, fmt.Errorf("[newMeiliSearchStore] update searchable attributes failed: %w", err)
	} else {
		logrus.Info("[newMeiliSearchStore] succesfuly update searchable attribute")
	}

	// UpdateSortableAttributes
	sortableAttributes := albumsearch.GetSortableAttributes()
	resp, err = albumIndex.UpdateSortableAttributes(&sortableAttributes)
	if err != nil {
		return nil, fmt.Errorf("[newMeiliSearchStore] update searchable attributes failed: %w", err)
	}
	task, err = cli.WaitForTask(resp.TaskUID, 0)
	if err == nil && task.Status == meilisearch.TaskStatusFailed {
		err = errors.New(task.Error.Message)
	}
	if err != nil {
		return nil, fmt.Errorf("[newMeiliSearchStore] update searchable attributes failed: %w", err)
	} else {
		logrus.Info("[newMeiliSearchStore] succesfuly update searchable attribute")
	}

	return &meiliSearchStore{
		cli:   cli,
		index: albumIndex,
	}, nil
}

type searchGetAlbumIDsOnlyResponse struct {
	Hits []struct {
		ID uint `json:"id"`
	} `json:"hits"`
}

func (s *meiliSearchStore) SearchGetIDsOrderByCreatedAt(ctx context.Context, query string, userID string, page, perPage int64, sortBy string) (albumIDs []uint, err error) {
	return s.searchGetIds(ctx, query, userID, page, perPage, albumsearch.ORDER_BY_CREATED_AT, sortBy)
}
func (s *meiliSearchStore) SearchGetIDsOrderByUploadedAt(ctx context.Context, query string, userID string, page, perPage int64, sortBy string) (albumIDs []uint, err error) {
	return s.searchGetIds(ctx, query, userID, page, perPage, albumsearch.ORDER_BY_UPLOADED_AT, sortBy)
}

func (s *meiliSearchStore) searchGetIds(ctx context.Context, query, userID string, page, perPage int64, orderBy string, sort string) (albumIDs []uint, err error) {
	sortArg := fmt.Sprintf("%s:%s", orderBy, sort)
	var rawResponse *json.RawMessage
	rawResponse, err = s.index.SearchRawWithContext(ctx, query, &meilisearch.SearchRequest{
		Sort:                 []string{sortArg},
		Page:                 page,
		HitsPerPage:          perPage,
		Filter:               fmt.Sprintf("%s = '%s'", albumsearch.USER_ID_KEY, userID),
		AttributesToRetrieve: []string{"id"},
	})
	if err != nil {
		return
	}
	var response searchGetAlbumIDsOnlyResponse
	err = json.Unmarshal(*rawResponse, &response)
	if err != nil {
		return
	}
	albumIDs = make([]uint, len(response.Hits))
	for i, hit := range response.Hits {
		albumIDs[i] = hit.ID
	}
	return
}

func (s *meiliSearchStore) Insert(ctx context.Context, documents []albumsearch.Index) (int64, error) {
	taskInfo, err := s.index.AddDocuments(documents, albumsearch.PRIMARY_KEY)
	if err != nil {
		return 0, err
	}
	return taskInfo.TaskUID, nil
}

func (s *meiliSearchStore) Delete(ctx context.Context, ids []string) (int64, error) {
	taskInfo, err := s.index.DeleteDocuments(ids)
	if err != nil {
		return 0, err
	}
	return taskInfo.TaskUID, nil
}

func (s *meiliSearchStore) MonitorTask(taskID int64, taskName string) error {
	task, err := s.cli.WaitForTask(taskID, 0)
	if err == nil && task.Status == meilisearch.TaskStatusFailed {
		err = errors.New(task.Error.Message)
	}
	logger := logrus.WithFields(logrus.Fields{
		"task_id": taskID,
		"index":   albumsearch.GetIndexName(),
	})
	if err != nil {
		logger.Errorf("[monitorTask] %s task failed: %v", taskName, err)
	}
	logger.Infof("[monitorTask] %s task completed", taskName)
	return err
}

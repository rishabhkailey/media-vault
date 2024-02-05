package mediasearchimpl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	mediasearch "github.com/rishabhkailey/media-service/internal/services/mediaSearch"
	"github.com/sirupsen/logrus"
)

type store interface {
	// media(2nd argument) pointer because gorm adds the missing info like ID, create_at to the pointer it self.
	Insert(context.Context, []mediasearch.Model) (int64, error)
	Delete(context.Context, []string) (int64, error)
	Search(context.Context, mediasearch.MediaSearchQuery) ([]mediasearch.Model, error)
	SearchGetMediaIDsOnly(context.Context, mediasearch.MediaSearchQuery) ([]uint, error)
	MonitorTask(int64, string) error
}

type meiliSearchStore struct {
	cli   *meilisearch.Client
	index *meilisearch.Index
}

var _ store = (*meiliSearchStore)(nil)

func newMeiliSearchStore(cli *meilisearch.Client) (*meiliSearchStore, error) {
	mediaIndex := cli.Index("media")

	// UpdateSearchableAttributes
	resp, err := mediaIndex.UpdateSearchableAttributes(&mediasearch.SearchableAttributes)
	if err != nil {
		return nil, fmt.Errorf("[newMeiliSearchStore] update searchable attributes failed: %w", err)
	}
	task, err := cli.WaitForTask(resp.TaskUID)
	if err == nil && task.Status == meilisearch.TaskStatusFailed {
		err = errors.New(task.Error.Message)
	}
	if err != nil {
		return nil, fmt.Errorf("[newMeiliSearchStore] update searchable attributes failed: %w", err)
	} else {
		logrus.Info("[newMeiliSearchStore] succesfuly update searchable attribute")
	}

	// UpdateFilterableAttributes
	resp, err = mediaIndex.UpdateFilterableAttributes(&mediasearch.FilterableAttributes)
	if err != nil {
		return nil, fmt.Errorf("[newMeiliSearchStore] update searchable attributes failed: %w", err)
	}
	task, err = cli.WaitForTask(resp.TaskUID)
	if err == nil && task.Status == meilisearch.TaskStatusFailed {
		err = errors.New(task.Error.Message)
	}
	if err != nil {
		return nil, fmt.Errorf("[newMeiliSearchStore] update searchable attributes failed: %w", err)
	} else {
		logrus.Info("[newMeiliSearchStore] succesfuly update searchable attribute")
	}

	// UpdateSortableAttributes
	resp, err = mediaIndex.UpdateSortableAttributes(&mediasearch.SortableAttributes)
	if err != nil {
		return nil, fmt.Errorf("[newMeiliSearchStore] update searchable attributes failed: %w", err)
	}
	task, err = cli.WaitForTask(resp.TaskUID)
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
		index: mediaIndex,
	}, nil
}

func (s *meiliSearchStore) Insert(ctx context.Context, documents []mediasearch.Model) (int64, error) {
	taskInfo, err := s.index.AddDocuments(documents, mediasearch.PRIMARY_KEY)
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
	task, err := s.cli.WaitForTask(taskID)
	if err == nil && task.Status == meilisearch.TaskStatusFailed {
		err = errors.New(task.Error.Message)
	}
	if err != nil {
		logrus.Errorf("[monitorTask] %s task failed: %v", taskName, err)
	}
	logrus.Infof("[monitorTask] %s task completed", taskName)
	return err
}

type searchResponse struct {
	Hits []mediasearch.Model `json:"hits"`
}

func (s *meiliSearchStore) Search(ctx context.Context, query mediasearch.MediaSearchQuery) (results []mediasearch.Model, err error) {
	orderByAttribute := mediasearch.OrderAttributesMapping[query.OrderBy]
	sort := fmt.Sprintf("%s:%s", orderByAttribute, query.Sort)
	var rawResponse *json.RawMessage
	rawResponse, err = s.index.SearchRaw(query.Query, &meilisearch.SearchRequest{
		Sort:        []string{sort},
		Page:        query.Page,
		HitsPerPage: query.PerPage,
		Filter:      fmt.Sprintf("%s = '%s'", mediasearch.USER_ID_KEY, query.UserID),
	})
	if err != nil {
		return
	}
	var response searchResponse
	err = json.Unmarshal(*rawResponse, &response)
	if err == nil {
		results = response.Hits
	}
	return
}

type searchGetMediaIDsOnlyResponse struct {
	Hits []struct {
		MediaID uint `json:"media_id"`
	} `json:"hits"`
}

func (s *meiliSearchStore) SearchGetMediaIDsOnly(ctx context.Context, query mediasearch.MediaSearchQuery) (mediaIDs []uint, err error) {
	orderByAttribute := mediasearch.OrderAttributesMapping[query.OrderBy]
	sort := fmt.Sprintf("%s:%s", orderByAttribute, mediasearch.SortKeywordMapping[query.Sort])
	var rawResponse *json.RawMessage
	rawResponse, err = s.index.SearchRaw(query.Query, &meilisearch.SearchRequest{
		Sort:                 []string{sort},
		Page:                 query.Page,
		HitsPerPage:          query.PerPage,
		Filter:               fmt.Sprintf("%s = '%s'", mediasearch.USER_ID_KEY, query.UserID),
		AttributesToRetrieve: []string{"media_id"},
	})
	if err != nil {
		return
	}
	var response searchGetMediaIDsOnlyResponse
	err = json.Unmarshal(*rawResponse, &response)
	if err != nil {
		return
	}
	for _, hit := range response.Hits {
		mediaIDs = append(mediaIDs, hit.MediaID)
	}
	return
}

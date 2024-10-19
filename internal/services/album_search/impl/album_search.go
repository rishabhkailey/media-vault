package albumsearchimpl

import (
	"context"
	"fmt"
	"strconv"

	"github.com/meilisearch/meilisearch-go"
	albumsearch "github.com/rishabhkailey/media-vault/internal/services/album_search"
	albumsearchstore "github.com/rishabhkailey/media-vault/internal/store/album_search"
	albumsearchimpl "github.com/rishabhkailey/media-vault/internal/store/album_search/impl"
	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
)

type Service struct {
	store albumsearchstore.Store
}

func NewService(cli meilisearch.ServiceManager) (albumsearch.Service, error) {
	store, err := albumsearchimpl.NewMeiliSearchStore(cli)
	if err != nil {
		return nil, fmt.Errorf("[NewService]: failed to create store: %w", err)
	}

	return &Service{
		store: store,
	}, nil
}

func (s *Service) GetAlbumIDs(ctx context.Context, args albumsearch.AlbumSearchQuery) ([]uint, error) {
	switch args.OrderBy {
	case "created_at":
		return s.store.SearchGetIDsOrderByCreatedAt(ctx, args.Query, args.UserID, args.Page, args.PerPage, args.Sort)
	case "uploaded_at":
		return s.store.SearchGetIDsOrderByUploadedAt(ctx, args.Query, args.UserID, args.Page, args.PerPage, args.Sort)
	default:
		return []uint{}, fmt.Errorf("[GetAlbumIDs]: invalid order by %s", args.OrderBy)
	}
}

func AlbumToMeiliAlbumSearchIndex(albums []storemodels.AlbumModel, userID string) (meiliSearchMediaList []albumsearchstore.Index) {
	for _, album := range albums {
		meiliSearchMediaList = append(meiliSearchMediaList, albumsearchstore.Index{
			Name:   album.Name,
			UserID: userID,
			ID:     album.ID,
		})
	}
	return
}

func (s *Service) CreateOne(ctx context.Context, album storemodels.AlbumModel, userId string) (taskID int64, err error) {
	taskID, err = s.store.Insert(ctx, AlbumToMeiliAlbumSearchIndex([]storemodels.AlbumModel{album}, userId))
	go s.store.MonitorTask(taskID, fmt.Sprintf("[MediaSearch] add document primaryKey=%v", album.ID))
	return
}

func (s *Service) CreateMany(ctx context.Context, albums []storemodels.AlbumModel, userId string) (taskID int64, err error) {
	taskID, err = s.store.Insert(ctx, AlbumToMeiliAlbumSearchIndex(albums, userId))
	ids := ""
	for _, mediaSearchData := range albums {
		ids += fmt.Sprintf(" %d", mediaSearchData.ID)
	}
	go s.store.MonitorTask(taskID, fmt.Sprintf("[MediaSearch] add document IDs=%v", ids))
	return
}

func (s *Service) DeleteOne(ctx context.Context, id int64) (taskID int64, err error) {
	taskID, err = s.store.Delete(ctx, []string{strconv.FormatUint(uint64(id), 10)})
	go s.store.MonitorTask(taskID, fmt.Sprintf("[MediaSearch] delete document primaryKey=%v", id))
	return
}

func (s *Service) DeleteMany(ctx context.Context, ids []int64) (taskID int64, err error) {
	idsStr := make([]string, 0, len(ids))
	for i, id := range ids {
		idsStr[i] = strconv.FormatUint(uint64(id), 10)
	}
	taskID, err = s.store.Delete(ctx, idsStr)
	go s.store.MonitorTask(taskID, fmt.Sprintf("[MediaSearch] delete documents primaryKey=%v", idsStr))
	return
}

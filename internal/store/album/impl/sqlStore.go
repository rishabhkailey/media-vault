package albumimpl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgconn"
	"github.com/rishabhkailey/media-vault/internal/constants"
	"github.com/rishabhkailey/media-vault/internal/store/album"
	storemodels "github.com/rishabhkailey/media-vault/internal/store/models"
	"github.com/rishabhkailey/media-vault/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type sqlStore struct {
	db    *gorm.DB
	cache *redis.Client
}

func (s *sqlStore) WithTransaction(tx *gorm.DB) album.Store {
	return &sqlStore{
		db:    tx,
		cache: s.cache,
	}
}

var _ album.Store = (*sqlStore)(nil)

func NewSqlStore(db *gorm.DB, cache *redis.Client) (album.Store, error) {
	if err := db.Migrator().AutoMigrate(&storemodels.AlbumModel{}); err != nil {
		return nil, err
	}
	if err := db.Migrator().AutoMigrate(&storemodels.AlbumMediaBindingsModel{}); err != nil {
		return nil, err
	}
	if err := db.Migrator().AutoMigrate(&storemodels.UserAlbumBindingsModel{}); err != nil {
		return nil, err
	}
	return &sqlStore{
		db:    db,
		cache: cache,
	}, nil
}

func (s *sqlStore) InsertAlbum(ctx context.Context, albumName string, thumbnailUrl string) (storemodels.AlbumModel, error) {
	album := storemodels.AlbumModel{
		Name:         albumName,
		ThumbnailUrl: thumbnailUrl,
	}
	err := s.db.WithContext(ctx).Create(&album).Error
	return album, err
}

func (s *sqlStore) InsertUserAlbumBindings(ctx context.Context, userID string, albumID uint) (uint, error) {
	album := storemodels.UserAlbumBindingsModel{
		UserID:  userID,
		AlbumID: albumID,
	}
	err := s.db.WithContext(ctx).Create(&album).Error
	return album.ID, err
}

func (s *sqlStore) GetByUserId(
	ctx context.Context,
	userID string,
	orderBy album.AlbumOrderBy,
	sort album.Sort,
	limit int,
	offset int,
) (albums []storemodels.AlbumModel, err error) {
	db := s.db.WithContext(ctx)
	// todo join?
	// cost of join
	// cost of join > cost of fetching all media_ids of the album and then filter
	albumByUserIDQuery := db.Model(&storemodels.UserAlbumBindingsModel{}).
		Select("album_id").
		Where("user_id = ?", userID)
	err = db.Model(&storemodels.AlbumModel{}).
		Where("id IN (?)", albumByUserIDQuery).
		Limit(limit).
		Order(fmt.Sprintf("%s %s", orderBy, sort)).
		Offset(offset).
		Find(&albums).Error
	return
}

func (s *sqlStore) GetAlbumsByUserIdOrderByCreationAt(
	ctx context.Context,
	userID string,
	orderBy album.AlbumOrderBy,
	sort album.Sort,
	lastAlbumID *uint,
	limit int,
) (albums []storemodels.AlbumModel, err error) {
	db := s.db.WithContext(ctx)
	albumsByUserIDQuery := db.Joins("Album").Model(&storemodels.UserAlbumBindingsModel{})

	if lastAlbumID != nil {
		var lastAlbumDate time.Time
		{
			err = db.Model(&storemodels.AlbumModel{}).
				Select("created_at").
				Where("id = @lastAlbumID", sql.Named("lastAlbumID", lastAlbumID)).
				First(&lastAlbumDate).Error
			if err != nil {
				return
			}
		}

		switch sort {
		case album.Ascending:
			{
				albumsByUserIDQuery = albumsByUserIDQuery.Where(`
					user_id = @userID
					AND (
						(
							("Album"."created_at" = @lastAlbumDate) AND ("Album"."id" < @lastAlbumID)
						) OR (
							("Album"."created_at" > @lastAlbumDate)
						)
						)`,
					sql.Named("userID", userID),
					sql.Named("lastAlbumDate", lastAlbumDate),
					sql.Named("lastAlbumID", lastAlbumID),
				)
			}
		default:
			{
				albumsByUserIDQuery = albumsByUserIDQuery.Where(`
					user_id = @userID 
					AND (
						(
							("Album"."created_at" = @lastAlbumDate) AND ("Album"."id" < @lastAlbumID)
						) OR (
							("Album"."created_at" < @lastAlbumDate)
						)
						)`,
					sql.Named("userID", userID),
					sql.Named("lastAlbumDate", lastAlbumDate),
					sql.Named("lastAlbumID", lastAlbumID),
				)
			}
		}
	} else {
		albumsByUserIDQuery = albumsByUserIDQuery.Where(
			`user_id = @userID `,
			sql.Named("userID", userID),
		)
	}
	queryOrderBy := fmt.Sprintf(`"Album"."created_at" %s, "Album"."id" desc`, sort)
	var userAlbumBindings []storemodels.UserAlbumBindingsModel
	err = albumsByUserIDQuery.
		Order(queryOrderBy).
		Limit(limit).
		Find(&userAlbumBindings).Error
	if err != nil {
		return
	}
	for _, userAlbumBinding := range userAlbumBindings {
		albums = append(albums, userAlbumBinding.Album)
	}
	return
}

func (s *sqlStore) GetAlbumsByUserIdOrderByUpdatedAt(
	ctx context.Context,
	userID string,
	orderBy album.AlbumOrderBy,
	sort album.Sort,
	lastAlbumID *uint,
	limit int,
) (albums []storemodels.AlbumModel, err error) {
	db := s.db.WithContext(ctx)
	albumsByUserIDQuery := db.Joins("Album").
		Model(&storemodels.UserAlbumBindingsModel{})

	if lastAlbumID != nil {
		var lastAlbumDate time.Time
		{
			err = db.Model(&storemodels.AlbumModel{}).
				Select("updated_at").
				Where("id = @lastAlbumID", sql.Named("lastAlbumID", lastAlbumID)).
				First(&lastAlbumDate).Error
			if err != nil {
				return
			}
		}

		switch sort {
		case album.Ascending:
			{
				albumsByUserIDQuery = albumsByUserIDQuery.Where(`
					user_id = @userID
					AND (
						(
							("Album"."updated_at" = @lastAlbumDate) AND ("Album"."id" < @lastAlbumID)
						) OR (
							("Album"."updated_at" > @lastAlbumDate)
						)
						)`,
					sql.Named("userID", userID),
					sql.Named("lastAlbumDate", lastAlbumDate),
					sql.Named("lastAlbumID", lastAlbumID),
				)
			}
		default:
			{
				albumsByUserIDQuery = albumsByUserIDQuery.Where(`
					user_id = @userID 
					AND (
						(
							("Album"."updated_at" = @lastAlbumDate) AND ("Album"."id" < @lastAlbumID)
						) OR (
							("Album"."updated_at" < @lastAlbumDate)
						)
						)`,
					sql.Named("userID", userID),
					sql.Named("lastAlbumDate", lastAlbumDate),
					sql.Named("lastAlbumID", lastAlbumID),
				)
			}
		}
	} else {
		albumsByUserIDQuery = albumsByUserIDQuery.Where(
			`user_id = @userID `,
			sql.Named("userID", userID),
		)
	}
	queryOrderBy := fmt.Sprintf(`"Album"."updated_at" %s, "Album"."id" desc`, sort)
	var userAlbumBindings []storemodels.UserAlbumBindingsModel
	err = albumsByUserIDQuery.
		Order(queryOrderBy).
		Limit(limit).
		Find(&userAlbumBindings).Error
	if err != nil {
		return
	}
	for _, userAlbumBinding := range userAlbumBindings {
		albums = append(albums, userAlbumBinding.Album)
	}
	return
}

func (s *sqlStore) GetByID(ctx context.Context, albumID uint) (result storemodels.AlbumModel, err error) {
	err = s.db.WithContext(ctx).Model(&storemodels.AlbumModel{}).First(&result, albumID).Error
	return
}

func (s *sqlStore) GetMediaByAlbumIdOrderByDate(ctx context.Context,
	albumID uint,
	lastMediaID *uint,
	sort album.Sort,
	limit int,
) (albumMediaBindings []storemodels.AlbumMediaBindingsModel, err error) {

	db := s.db.WithContext(ctx)

	// table name = media, metadata alias = Metadata
	query := db.
		Preload("Media.Metadata").
		Joins(
			`LEFT JOIN "media" on "album_media_bindings"."media_id" = "media"."id"
		AND "album_media_bindings"."deleted_at" IS NULL`,
		).
		Joins(
			`LEFT JOIN "media_metadata" ON "media"."metadata_id" = "media_metadata"."id"
			AND "media_metadata"."deleted_at" IS NULL`,
		).
		Model(&storemodels.AlbumMediaBindingsModel{})

	if lastMediaID != nil {

		var lastMediaDate time.Time
		{
			// it is adding extra select columns not sure why
			// err = db.
			// 	Joins("Metadata").Model(&media.Media{}).
			// 	Select(`"Metadata"."date"`).
			// 	Where(
			// 		`"media"."id" = @lastMediaID`,
			// 		sql.Named("lastMediaID", lastMediaID),
			// 	).Find(&lastMediaDate).Error
			err = db.Table(constants.MEDIA_METADATA_TABLE).
				Select("date").
				Where("id = (@lastMediaMetadataID)",
					sql.Named(
						"lastMediaMetadataID",
						db.Model(&storemodels.MediaModel{}).
							Select("metadata_id").
							Where("id = ?", lastMediaID),
					),
				).First(&lastMediaDate).Error
			if err != nil {
				return
			}
		}

		switch sort {
		case album.Ascending:
			{
				query = query.Where(`
					"album_media_bindings"."album_id" = @albumID
					AND (
						(
							("media_metadata"."date" = @lastMediaDate) AND ("media"."id" < @lastMediaID)
						) OR (
							("media_metadata"."date" > @lastMediaDate)
						)
					)`,
					sql.Named("albumID", albumID),
					sql.Named("lastMediaDate", lastMediaDate),
					sql.Named("lastMediaID", lastMediaID),
				)
			}
		default:
			{
				query = query.Where(`
					"album_media_bindings"."album_id" = @albumID
					AND (
						(
							("media_metadata"."date" = @lastMediaDate) AND ("media"."id" < @lastMediaID)
						) OR (
							("media_metadata"."date" < @lastMediaDate)
						)
					)`,
					sql.Named("albumID", albumID),
					sql.Named("lastMediaDate", lastMediaDate),
					sql.Named("lastMediaID", lastMediaID),
				)
			}
		}
	} else {
		query = query.Where(
			`"album_media_bindings"."album_id" = @albumID`,
			sql.Named("albumID", albumID),
		)
	}

	queryOrderBy := fmt.Sprintf(`"media_metadata"."date" %s, "media"."id" desc`, sort)
	err = query.
		Order(queryOrderBy).
		Limit(limit).
		Find(&albumMediaBindings).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []storemodels.AlbumMediaBindingsModel{}, nil
	}
	return
}

func (s *sqlStore) GetMediaByAlbumIdOrderByUploadDate(ctx context.Context,
	albumID uint,
	lastMediaID *uint,
	sort album.Sort,
	limit int,
) (albumMediaBindings []storemodels.AlbumMediaBindingsModel, err error) {
	db := s.db.WithContext(ctx)

	query := db.
		Preload("Media.Metadata").
		Joins(
			`LEFT JOIN "media" on "album_media_bindings"."media_id" = "media"."id"
		AND "album_media_bindings"."deleted_at" IS NULL`,
		).
		Joins(
			`LEFT JOIN "media_metadata" ON "media"."metadata_id" = "media_metadata"."id"
			AND "media_metadata"."deleted_at" IS NULL`,
		).
		Model(&storemodels.AlbumMediaBindingsModel{})

	if lastMediaID != nil {
		var lastUploadDate time.Time
		{
			err = db.Model(&storemodels.MediaModel{}).
				Select("created_at").
				Where(
					"id = @lastMediaID",
					sql.Named("lastMediaID", lastMediaID),
				).First(&lastUploadDate).Error
			if err != nil {
				return
			}
		}

		switch sort {
		case album.Ascending:
			{
				// media whose created date > last media date (sort by created_at asc)
				// media whose created date = last media date and media id < last media id (media id is sorted by desc order)
				query = query.Where(`
					"album_media_bindings"."album_id" = @albumID
					AND (
						(
							("media_metadata"."created_at" = @lastUploadDate) AND ("media"."id" < @lastMediaID)
						) OR (
							("media_metadata"."created_at" > @lastUploadDate)
						)
					)`,
					sql.Named("albumID", albumID),
					sql.Named("lastUploadDate", lastUploadDate),
					sql.Named("lastMediaID", lastMediaID),
				)
			}
		default:
			{
				// media whose created date < last media date (sort by created_at desc)
				// media whose created date = last media date and media id < last media id (media id is sorted by desc order)
				query = query.Where(`
					"album_media_bindings"."album_id" = @albumID
					AND (
						(
							("media_metadata"."created_at" = @lastUploadDate) AND ("media"."id" < @lastMediaID)
						) OR (
							("media_metadata"."created_at" < @lastUploadDate)
						)
					)`,
					sql.Named("albumID", albumID),
					sql.Named("lastUploadDate", lastUploadDate),
					sql.Named("lastMediaID", lastMediaID),
				)
			}
		}
	} else {
		query = query.Where(
			`"album_media_bindings"."album_id" = @albumID`,
			sql.Named("albumID", albumID),
		)
	}

	queryOrderBy := fmt.Sprintf(`"media_metadata"."created_at" %s, "media"."id" desc`, sort)
	err = query.
		Order(queryOrderBy).
		Limit(limit).
		Find(&albumMediaBindings).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []storemodels.AlbumMediaBindingsModel{}, nil
	}
	return
}

func (s *sqlStore) GetMediaByAlbumIdOrderByAddedDate(
	ctx context.Context,
	albumID uint,
	lastMediaID *uint,
	sort album.Sort,
	limit int,
) (albumMediaBindings []storemodels.AlbumMediaBindingsModel, err error) {
	db := s.db.WithContext(ctx)
	// SELECT media_id FROM album_media_bindings
	query := db.
		Preload("Media.Metadata").
		Joins(
			`LEFT JOIN "media" on "album_media_bindings"."media_id" = "media"."id"
		AND "album_media_bindings"."deleted_at" IS NULL`,
		).
		Joins(
			`LEFT JOIN "media_metadata" ON "media"."metadata_id" = "media_metadata"."id"
			AND "media_metadata"."deleted_at" IS NULL`,
		).
		Model(&storemodels.AlbumMediaBindingsModel{})

	if lastMediaID != nil {
		// SELECT created_at FROM album_media_bindings where album_id = ? AND media_id = ?
		// lastDate not being used because we are not returning the added to album date for a media in response
		var lastAddedAtDate time.Time
		{

			err = db.Model(&storemodels.AlbumMediaBindingsModel{}).
				Select("created_at").
				Where(
					`"media_id" = @lastMediaID AND album_id = @albumID`,
					sql.Named("albumID", albumID),
					sql.Named("lastMediaID", lastMediaID),
				).
				First(&lastAddedAtDate).Error
			if err != nil {
				return
			}
		}
		switch sort {
		case album.Ascending:
			{
				query = query.Where(`"album_id" = @albumID
			AND (
				(
					("album_media_bindings"."created_at" = @lastAddedAtDate) AND ("album_media_bindings"."media_id" < @lastMediaID)
				) OR (
					("album_media_bindings"."created_at" > @lastAddedAtDate)
				)
			)
			`,
					sql.Named("albumID", albumID),
					sql.Named("lastAddedAtDate", lastAddedAtDate),
					sql.Named("lastMediaID", lastMediaID),
				)
			}
		default:
			{
				query = query.Where(`"album_media_bindings"."album_media_bindings""album_id" = @albumID
			AND (
				(
					("album_media_bindings"."album_media_bindings""created_at" = @lastAddedAtDate) AND ("album_media_bindings"."media_id" < @lastMediaID)
				) OR (
					("album_media_bindings"."album_media_bindings""created_at" < @lastAddedAtDate)
				)
			)
			`,
					sql.Named("albumID", albumID),
					sql.Named("lastAddedAtDate", lastAddedAtDate),
					sql.Named("lastMediaID", lastMediaID),
				)
			}
		}
	} else {
		query.Where(
			`"album_media_bindings"."album_id" = @albumID`,
			sql.Named("albumID", albumID),
		)
	}

	queryOrderBy := fmt.Sprintf(`"album_media_bindings"."created_at" %s, "album_media_bindings"."media_id" desc`, sort)

	err = query.
		Order(queryOrderBy).
		Limit(limit).
		Find(&albumMediaBindings).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []storemodels.AlbumMediaBindingsModel{}, nil
	}

	return
}

func (s *sqlStore) CheckAlbumBelongsToUser(ctx context.Context, userID string, albumID uint) (ok bool, err error) {
	var userAlbumBinding storemodels.UserAlbumBindingsModel
	err = s.db.WithContext(ctx).First(&userAlbumBinding, "album_id = ?", albumID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return userAlbumBinding.UserID == userID, nil
}

func (s *sqlStore) AddMediaInAlbum(ctx context.Context, albumID uint, mediaIDs []uint) (newMediaIDs []uint, err error) {
	alreadyExist, err := s.FilterMediaIDsByAlbumID(ctx, albumID, mediaIDs)
	if err != nil {
		return
	}
	newMediaIDs = make([]uint, 0)
	var albumMediaBindings []storemodels.AlbumMediaBindingsModel
	for _, mediaID := range mediaIDs {
		if utils.Contains(alreadyExist, mediaID) {
			continue
		}
		newMediaIDs = append(newMediaIDs, mediaID)
		albumMediaBindings = append(albumMediaBindings, storemodels.AlbumMediaBindingsModel{
			AlbumID: albumID,
			MediaID: mediaID,
		})
	}
	if len(albumMediaBindings) == 0 {
		return
	}
	err = s.db.WithContext(ctx).CreateInBatches(albumMediaBindings, 100).Error
	if isUniqueConstraintError(err) {
		logrus.Warnf("[AddMediaInAlbum] duplicate constraint ignoring some rows: %v", err)
		err = nil
	}
	if err != nil {
		return
	}
	if _, err := s.UpdateAlbumMediaCount(ctx, albumID, len(newMediaIDs)); err != nil {
		logrus.Warnf("[AddMediaInAlbum] album updatedAt failed: %v", err)
	}
	return
}

func (s *sqlStore) RemoveMediaFromAlbum(ctx context.Context, albumID uint, mediaIDs []uint) (removedMediaIDs []uint, err error) {
	mediaIDs, err = s.FilterMediaIDsByAlbumID(ctx, albumID, mediaIDs)
	if err != nil {
		return
	}
	db := s.db.WithContext(ctx).Begin()
	err = db.Unscoped().Delete(&storemodels.AlbumMediaBindingsModel{}, "album_id = ? AND media_id IN ?", albumID, mediaIDs).Error
	if err != nil {
		db.Rollback()
		return
	}
	db.Commit()
	removedMediaIDs = mediaIDs
	if _, err := s.UpdateAlbumMediaCount(ctx, albumID, len(removedMediaIDs)*-1); err != nil {
		logrus.Warnf("[AddMediaInAlbum] album updatedAt failed: %v", err)
	}
	return
}

func (s *sqlStore) FilterMediaIDsByAlbumID(ctx context.Context, albumID uint, mediaIDs []uint) (result []uint, err error) {
	err = s.db.Model(&storemodels.AlbumMediaBindingsModel{}).Where("media_id IN ? and album_id = ?", mediaIDs, albumID).Pluck("media_id", &result).Error
	return
}

func (s *sqlStore) DeleteAlbum(ctx context.Context, albumID uint) error {
	tx := s.db.Begin().WithContext(ctx)

	if err := tx.Unscoped().Delete(&storemodels.UserAlbumBindingsModel{}, "album_id = ?", albumID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Delete(&storemodels.AlbumMediaBindingsModel{}, "album_id = ?", albumID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Delete(&storemodels.AlbumModel{
		Model: gorm.Model{
			ID: albumID,
		},
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *sqlStore) UpdateAlbum(ctx context.Context, albumID uint, name *string, thumbnailUrl *string) (updatedAlbum storemodels.AlbumModel, err error) {
	err = s.db.WithContext(ctx).First(&updatedAlbum, albumID).Error
	if err != nil {
		return
	}
	if name != nil {
		updatedAlbum.Name = *name
	}
	if thumbnailUrl != nil {
		updatedAlbum.ThumbnailUrl = *thumbnailUrl
	}
	err = s.db.WithContext(ctx).Save(&updatedAlbum).Error
	return
}

func (s *sqlStore) UpdateAlbumMediaCount(ctx context.Context, albumID uint, change int) (updatedAlbum storemodels.AlbumModel, err error) {
	// album.ID = albumID
	// album.UpdatedAt = updatedAt

	err = s.db.WithContext(ctx).First(&updatedAlbum, albumID).Error
	if err != nil {
		return
	}
	updatedAlbum.MediaCount += change
	if updatedAlbum.MediaCount < 0 {
		updatedAlbum.MediaCount = 0
	}
	err = s.db.WithContext(ctx).Save(&updatedAlbum).Error
	return
}

func isUniqueConstraintError(err error) bool {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		if pgError.Code == "23505" {
			return true
		}
	}
	return false
}

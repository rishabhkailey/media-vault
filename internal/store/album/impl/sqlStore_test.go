package albumimpl

import (
	"context"
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rishabhkailey/media-service/internal/store/album"
	storemodels "github.com/rishabhkailey/media-service/internal/store/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// mostly used for null
type AnyValue struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyValue) Match(v driver.Value) bool {
	return true
}

func NewMockSqlStore() (store album.Store, sqlMock sqlmock.Sqlmock, err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		"test", "test", "test", "test", 5432)
	conn, sqlMock, err := sqlmock.NewWithDSN(dsn)
	if err != nil {
		return
	}
	sqlLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Nanosecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       conn,
		DriverName: "postgres",
	}), &gorm.Config{
		Logger: sqlLogger,
	})
	if err != nil {
		return
	}
	// can not use newSqlStore mehtod as that will also try to run query to create missing tables and indexes
	store = &sqlStore{
		db:    db,
		cache: nil,
	}
	if err != nil {
		return
	}
	return
}

func TestAlbumInsert(t *testing.T) {
	store, sqlMock, err := NewMockSqlStore()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	album := storemodels.AlbumModel{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{
				Valid: false,
			},
		},
		Name:         "test_album",
		ThumbnailUrl: "/v1/thumbnail/test_album",
	}
	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO "albums" ("created_at","updated_at","deleted_at","name","thumbnail_url") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`,
		),
	).WithArgs(
		AnyTime{}, AnyTime{}, sqlmock.AnyArg(), album.Name, album.ThumbnailUrl,
	).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))
	sqlMock.ExpectCommit()
	if _, err := store.InsertAlbum(context.Background(), album.Name, album.ThumbnailUrl); err != nil {
		t.Fail()
		t.Error(err)
		return
	}
	err = sqlMock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

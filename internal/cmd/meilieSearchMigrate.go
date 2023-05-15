package cmd

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/meilisearch/meilisearch-go"
	mediasearch "github.com/rishabhkailey/media-service/internal/services/mediaSearch"
	mediasearchimpl "github.com/rishabhkailey/media-service/internal/services/mediaSearch/mediaSearchimpl"
	usermediabindings "github.com/rishabhkailey/media-service/internal/services/userMediaBindings"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func getMediaBatch(db *gorm.DB, offset int, limit int) (mediaList []usermediabindings.Model, err error) {
	// mediaIdQuery := db.Model(&services.UserMediaBinding{}).Select("media_id").Limit(limit).Offset(offset)
	// err = db.Joins("Metadata").Joins("Media").Model(&services.UserMediaBinding{}).Limit(limit).Offset(offset).Find(&mediaList).Error
	err = db.Preload("Media.Metadata").Model(&usermediabindings.Model{}).Limit(limit).Offset(offset).Find(&mediaList).Error
	return
}

func MeiliSearchMigrate(gormDB *gorm.DB, meiliSearch *meilisearch.Client, batchSize int) error {
	var total int64 = 0
	if err := gormDB.Model(&usermediabindings.Model{}).Count(&total).Error; err != nil {
		return fmt.Errorf("[MeiliSearchMigrate] count failed: %w", err)
	}
	var offset int64 = 0
	var batchNumber int = 0
	var waitGroup sync.WaitGroup
	mediaSearchService, err := mediasearchimpl.NewService(meiliSearch)
	if err != nil {
		return err
	}
	for offset < total {
		mediaList, err := getMediaBatch(gormDB, int(offset), batchSize)
		if err != nil {
			return fmt.Errorf("[MeiliSearchMigrate] batch number %v failed: %w", batchNumber, err)
		}
		// check after parsing as well both length should be same
		if len(mediaList) == 0 {
			logrus.Warnf("[MeiliSearchMigrate] batch %v with 0 size", batchNumber)
		}
		meiliSearchMediaList, err := mediasearch.UserMediaBindingToMeiliSearchMediaIndex(mediaList)
		if err != nil {
			return fmt.Errorf("[MeiliSearchMigrate] toMeiliSearchMediaIndex failed for batch %v: %w", batchNumber, err)
		}
		if len(meiliSearchMediaList) != len(mediaList) {
			return fmt.Errorf("[MeiliSearchMigrate] meiliSearchMediaList returned slice of different length for batch %v", batchNumber)
		}
		taskID, err := mediaSearchService.CreateMany(context.Background(), meiliSearchMediaList)
		if err != nil {
			return fmt.Errorf("[MeiliSearchMigrate] add documents failed for batch %v: %w", batchNumber, err)
		}
		logrus.Infof("[MeiliSearchMigrate] batch %d is proccessing", batchNumber)
		waitGroup.Add(1)
		go func(batchNumber int) {
			task, err := meiliSearch.WaitForTask(taskID)
			if err == nil && task.Status == meilisearch.TaskStatusFailed {
				err = errors.New(task.Error.Message)
			}
			if err != nil {
				logrus.Errorf("[MeiliSearchMigrate] batch %d failed: %w", batchNumber, err)
			} else {
				logrus.Infof("[MeiliSearchMigrate] batch %d %v", batchNumber, task.Status)
			}
			waitGroup.Done()
		}(batchNumber)
		offset += int64(batchSize)
		batchNumber++
	}
	waitGroup.Wait()
	logrus.Infof("[MeiliSearchMigrate] migration completed")
	return nil
}

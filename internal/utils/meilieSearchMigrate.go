package utils

import (
	"errors"
	"fmt"
	"sync"

	"github.com/meilisearch/meilisearch-go"
	"github.com/rishabhkailey/media-service/internal/db"
	"github.com/rishabhkailey/media-service/internal/db/services"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func getMediaBatch(db *gorm.DB, offset int, limit int) (mediaList []services.UserMediaBinding, err error) {
	// mediaIdQuery := db.Model(&services.UserMediaBinding{}).Select("media_id").Limit(limit).Offset(offset)
	// err = db.Joins("Metadata").Joins("Media").Model(&services.UserMediaBinding{}).Limit(limit).Offset(offset).Find(&mediaList).Error
	err = db.Preload("Media.Metadata").Model(&services.UserMediaBinding{}).Limit(limit).Offset(offset).Find(&mediaList).Error
	return
}

func MeiliSearchMigrate(gormDB *gorm.DB, meiliSearch *meilisearch.Client, batchSize int) error {
	var total int64 = 0
	if err := gormDB.Model(&services.UserMediaBinding{}).Count(&total).Error; err != nil {
		return fmt.Errorf("[MeiliSearchMigrate] count failed: %w", err)
	}
	var offset int64 = 0
	var batchNumber int = 0
	var waitGroup sync.WaitGroup
	mediaIndex := meiliSearch.Index("media")
	resp, err := mediaIndex.UpdateSearchableAttributes(&db.MeilieSearchMediaIndexSearchable)
	if err != nil {
		return fmt.Errorf("[MeiliSearchMigrate] update searchable attributes failed: %w", err)
	}
	// wait for UpdateSearchableAttributes
	{
		task, err := meiliSearch.WaitForTask(resp.TaskUID)
		if err == nil && task.Status == meilisearch.TaskStatusFailed {
			err = errors.New(task.Error.Message)
		}
		if err != nil {
			return fmt.Errorf("[MeiliSearchMigrate] update searchable attributes failed: %w", err)
		} else {
			logrus.Info("[MeiliSearchMigrate] succesfuly update searchable attribute")
		}
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
		meiliSearchMediaList, err := toMeiliSearchMediaIndex(mediaList)
		if err != nil {
			return fmt.Errorf("[MeiliSearchMigrate] toMeiliSearchMediaIndex failed for batch %v: %w", batchNumber, err)
		}
		if len(meiliSearchMediaList) != len(mediaList) {
			return fmt.Errorf("[MeiliSearchMigrate] meiliSearchMediaList returned slice of different length for batch %v", batchNumber)
		}
		taskInfo, err := mediaIndex.AddDocuments(meiliSearchMediaList, "media_id")
		if err != nil {
			return fmt.Errorf("[MeiliSearchMigrate] add documents failed for batch %v: %w", batchNumber, err)
		}
		logrus.Infof("[MeiliSearchMigrate] batch %d status is %v", batchNumber, taskInfo.Status)
		waitGroup.Add(1)
		go func(batchNumber int) {
			task, err := meiliSearch.WaitForTask(taskInfo.TaskUID)
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

func toMeiliSearchMediaIndex(userMediaBindingList []services.UserMediaBinding) (meiliSearchMediaList []db.MeiliSearchMediaIndex, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[toMeiliSearchMediaIndex] panic :%v", r)
		}
	}()
	for _, userMediaBinding := range userMediaBindingList {
		meiliSearchMediaList = append(meiliSearchMediaList, db.MeiliSearchMediaIndex{
			MediaID: userMediaBinding.Media.ID,
			UserID:  userMediaBinding.UserID,
			Metadata: db.MeiliSearchMediaMetadata{
				Name:      userMediaBinding.Media.Metadata.Name,
				Type:      userMediaBinding.Media.Metadata.Type,
				Timestamp: userMediaBinding.Media.Metadata.Date.Unix(),
				Date:      userMediaBinding.Media.Metadata.Date.Format("Monday January 2 2006 UTC"),
			},
		})
	}
	return
}

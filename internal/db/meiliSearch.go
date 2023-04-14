package db

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/rishabhkailey/media-service/internal/config"
)

type MeiliSearchMediaMetadata struct {
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
	Date      string `json:"date"`
	Type      string `json:"type"`
}

type MeiliSearchMediaIndex struct {
	MediaID  uint                     `json:"media_id"`
	UserID   string                   `json:"user_id"`
	Metadata MeiliSearchMediaMetadata `json:"metadata"`
}

var MeilieSearchMediaIndexSearchable = []string{"metadata"}

func NewMeiliSearchClient(config config.MeiliSearch) (*meilisearch.Client, error) {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   config.Host,
		APIKey: config.APIKey,
	})
	if _, err := client.Health(); err != nil {
		return nil, err
	}
	return client, nil
}

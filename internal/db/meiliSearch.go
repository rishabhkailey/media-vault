package db

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/rishabhkailey/media-service/internal/config"
)

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

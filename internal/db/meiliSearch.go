package db

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/rishabhkailey/media-vault/internal/config"
)

func NewMeiliSearchClient(config config.MeiliSearch) (meilisearch.ServiceManager, error) {
	client := meilisearch.New(config.Host, meilisearch.WithAPIKey(config.APIKey))
	if _, err := client.Health(); err != nil {
		return nil, err
	}
	return client, nil
}

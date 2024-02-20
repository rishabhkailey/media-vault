package config

import (
	"strings"

	"github.com/meilisearch/meilisearch-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RedisCacheConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Db       int
}

type ServerConfig struct {
	Host    string
	Port    int
	BaseURL string
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

type Session struct {
	Secret string
}

type WebUIConfig struct {
	Directory string
	BaseURL   string
}

type MinioConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Bucket   string
	// todo create struct for tl
	TLS              bool
	TLSSkipVerify    bool
	CustomRootCAPath string
}

type OIDC struct {
	MediaVault struct {
		ClientID     string
		ClientSecret string
	}
	SPA struct {
		ClientID string
	}
	DiscoveryEndpoint string
	URL               string
}

type MeiliSearch struct {
	meilisearch.ClientConfig
}

type Config struct {
	Cache RedisCacheConfig
	// JWT        *JWTConfig
	Server      ServerConfig
	Database    Database
	Session     Session
	WebUIConfig WebUIConfig
	MinioConfig MinioConfig
	OIDC        OIDC
	MeiliSearch MeiliSearch
}

// todo validation
func GetConfig() (*Config, error) {
	err := configInit()
	if err != nil {
		return nil, err
	}

	config := Config{
		Cache: RedisCacheConfig{
			Host:     viper.GetString("cache.redis.host"),
			Port:     viper.GetInt("cache.redis.port"),
			User:     viper.GetString("cache.redis.user"),
			Password: viper.GetString("cache.redis.password"),
			Db:       viper.GetInt("cache.redis.db"),
		},
		Server: ServerConfig{
			Host: viper.GetString("server.host"),
			Port: viper.GetInt("server.port"),
		},
		Database: Database{
			Host:     viper.GetString("database.postgres.host"),
			Port:     viper.GetInt("database.postgres.port"),
			User:     viper.GetString("database.postgres.user"),
			Password: viper.GetString("database.postgres.password"),
			Dbname:   viper.GetString("database.postgres.dbname"),
		},
		Session: Session{
			Secret: viper.GetString("session.secret"),
		},
		MinioConfig: MinioConfig{
			Host:             viper.GetString("minio.host"),
			Port:             viper.GetInt("minio.port"),
			User:             viper.GetString("minio.user"),
			Bucket:           viper.GetString("minio.bucket"),
			Password:         viper.GetString("minio.password"),
			TLS:              viper.GetBool("minio.tls.enabled"),
			TLSSkipVerify:    viper.GetBool("minio.tls.skipVerify"),
			CustomRootCAPath: viper.GetString("minio.tls.customRootCAPath"),
		},
		OIDC: OIDC{
			URL:               viper.GetString("oidc.url"),
			DiscoveryEndpoint: viper.GetString("oidc.discoveryEndpoint"),
			MediaVault: struct {
				ClientID     string
				ClientSecret string
			}{
				ClientID:     viper.GetString("oidc.mediaVault.client.id"),
				ClientSecret: viper.GetString("oidc.mediaVault.client.secret"),
			},
			SPA: struct{ ClientID string }{
				ClientID: viper.GetString("oidc.spa.client.id"),
			},
		},
		MeiliSearch: MeiliSearch{
			meilisearch.ClientConfig{
				Host:   viper.GetString("meiliSearch.host"),
				APIKey: viper.GetString("meiliSearch.APIKey"),
			},
		},
		WebUIConfig: WebUIConfig{
			// todo: add check for directory so user doesn't accidently expose sensitive files
			// check for the files only allow .css, .html, .js and other static files
			// also add a option to disable this check
			Directory: viper.GetString("webUI.directory"),
		},
	}
	return &config, nil
}

func configInit() error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MV")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// default config file path
	viper.SetConfigFile("configs/mediavault.yaml")
	// from cmd flag or env variable
	configFilePath := viper.GetString("config")
	if len(configFilePath) != 0 {
		viper.SetConfigFile(configFilePath)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logrus.Warnf("[configInit] unable to read config file: %v, application will try to read config from env variables", err)
		}
	}
	return nil
}

// todo validate function

package config

import (
	"fmt"
	"os"

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
	// todo create struct for tl
	TLS              bool
	TLSSkipVerify    bool
	CustomRootCAPath string
}

type Config struct {
	Cache RedisCacheConfig
	// JWT        *JWTConfig
	Server      ServerConfig
	Database    Database
	Session     Session
	WebUIConfig WebUIConfig
	MinioConfig MinioConfig
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
			Host:    viper.GetString("server.host"),
			Port:    viper.GetInt("server.port"),
			BaseURL: viper.GetString("server.baseURL"),
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
			Password:         viper.GetString("minio.password"),
			TLS:              viper.GetBool("minio.tls.enabled"),
			TLSSkipVerify:    viper.GetBool("minio.tls.skipVerify"),
			CustomRootCAPath: viper.GetString("minio.tls.customRootCAPath"),
		},
	}
	return &config, nil
}

func configInit() error {
	wd, _ := os.Getwd()
	_ = wd
	viper.SetEnvPrefix("AUTH_SERVICE")
	// default config file path
	viper.SetConfigFile("configs/authservice.yaml")
	// from cmd flag or env variable
	configFilePath := viper.GetString("config")
	if len(configFilePath) != 0 {
		viper.SetConfigFile(configFilePath)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("[GetConfig]: Config file not found: %w", err)
		}
		return fmt.Errorf("[GetConfig]: Config file found but unable to read: %w", err)
	}
	return nil
}

// validate function??

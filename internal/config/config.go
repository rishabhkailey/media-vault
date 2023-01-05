package config

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
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

// type JWTConfig struct {
// 	privateKey string
// 	publicKey  string
// }

type JWTSigningKey struct {
	Kid             string
	PrivateKey      *rsa.PrivateKey
	PublicKey       *rsa.PublicKey
	PrivateKeyBytes []byte
	PublicKeyBytes  []byte
}

// todo multiple signing keys for rotation
// type JWTSigningKeys struct {

// }

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

type Config struct {
	Cache RedisCacheConfig
	// JWT        *JWTConfig
	Server        ServerConfig
	Database      Database
	JWTSigningKey *JWTSigningKey
	Session       Session
	WebUIConfig   WebUIConfig
}

// todo validation
func GetConfig() (*Config, error) {
	err := configInit()
	if err != nil {
		return nil, err
	}

	signingKey, err := signingKeyFromBase64(viper.GetString("jwt.privateKey"), viper.GetString("jwt.publicKey"))
	if err != nil {
		return nil, err
	}

	webUIConfig := WebUIConfig{
		Directory: viper.GetString("webUI.directory"),
		BaseURL:   viper.GetString("webUI.baseURL"),
	}
	if len(webUIConfig.BaseURL) == 0 {
		webUIConfig.BaseURL = viper.GetString("server.baseURL")
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
		JWTSigningKey: signingKey,
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
		WebUIConfig: webUIConfig,
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

func signingKeyFromBase64(b64PrivateKey, b64PublicKey string) (*JWTSigningKey, error) {
	// todo add support for private key with password

	privateKeyBytes, err := base64.StdEncoding.DecodeString(viper.GetString("jwt.privateKey"))
	if err != nil {
		return nil, fmt.Errorf("error decoding base64 private key from config: %w", err)
	}
	publicKeyBytes, err := base64.StdEncoding.DecodeString(viper.GetString("jwt.publicKey"))
	if err != nil {
		return nil, fmt.Errorf("error decoding base64 public key from config: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid public key: %w", err)
	}
	return &JWTSigningKey{
		Kid:             "0",
		PublicKey:       publicKey,
		PrivateKey:      privateKey,
		PublicKeyBytes:  publicKeyBytes,
		PrivateKeyBytes: privateKeyBytes,
	}, nil
}

// validate function??

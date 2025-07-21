package config

import (
	"fmt"
	"marketplace/internal/models"
	"net"
	"os"

	"github.com/subosito/gotenv"
	"gopkg.in/yaml.v3"
)

type serverConfig struct {
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
	JWTSecretKey     string `yaml:"jwt_secret_key"`
	AESSecretKey     string `yaml:"aes_secret_key"`
	DefoultFeedLimit uint64 `yaml:"feed_limit"`
}

type tracerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type dbConfig struct {
	DbHost     string `yaml:"db_host"`
	DbPort     string `yaml:"db_port"`
	DbName     string `yaml:"db_name"`
	DbUser     string `yaml:"db_user"`
	DbPassword string `yaml:"db_password"`
	DbSSL      string `yaml:"db_sslmode"`
}

type MainConfig struct {
	ServerConfig serverConfig `yaml:"server"`
	TracerConfig tracerConfig `yaml:"tracer"`
	DbConfig     dbConfig     `yaml:"db"`
}

func InitMainConfig(configPath *string, isLocal *bool) (*MainConfig, error) {
	var cfg MainConfig

	if configPath == nil || isLocal == nil {
		return &cfg, models.ErrParseConfig
	}

	if _, err := os.Stat(*configPath); err != nil {
		return nil, err
	}
	rowConfig, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(rowConfig, &cfg)
	if err != nil {
		return nil, err
	}

	if *isLocal {
		err = gotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	cfg.ServerConfig.JWTSecretKey = os.Getenv("APP_JWT_SECRET_KEY")
	cfg.ServerConfig.AESSecretKey = os.Getenv("APP_AES_SECRET_KEY")

	cfg.DbConfig = dbConfig{
		DbHost:     os.Getenv("PG_HOST"),
		DbPort:     os.Getenv("PG_PORT"),
		DbName:     os.Getenv("PG_DATABASE_NAME"),
		DbUser:     os.Getenv("PG_USER"),
		DbPassword: os.Getenv("PG_PASSWORD"),
		DbSSL:      os.Getenv("PG_SSLMODE"),
	}

	return &cfg, nil
}

func (cfg *MainConfig) GetServerAddress() string {
	return net.JoinHostPort(cfg.ServerConfig.Host, cfg.ServerConfig.Port)
}

func (cfg *MainConfig) GetJWSKey() string {
	return cfg.ServerConfig.JWTSecretKey
}

func (cfg *MainConfig) GetAESKey() string {
	return cfg.ServerConfig.AESSecretKey
}

func (cfg *MainConfig) GetDefoultFeedLimit() uint64 {
	return cfg.ServerConfig.DefoultFeedLimit
}

func (cfg *MainConfig) GetTracerAddress() string {
	return net.JoinHostPort(cfg.TracerConfig.Host, cfg.TracerConfig.Port)
}

func (cfg *MainConfig) GetDbConfig() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.DbConfig.DbHost,
		cfg.DbConfig.DbPort,
		cfg.DbConfig.DbName,
		cfg.DbConfig.DbUser,
		cfg.DbConfig.DbPassword,
		cfg.DbConfig.DbSSL,
	)
}

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const defaultEnvironment = "test"

var Conf *TomlConfig

type TomlConfig struct {
	Common    CommonConfig    `toml:"COMMON"`
	MySQL     MySQLConfig     `toml:"MySQL"`
	MongoDB   MongoDBConfig   `toml:"MongoDB"`
	Redis     RedisConfig     `toml:"Redis"`
	JWT       JWTConfig       `toml:"JWT"`
	RateLimit RateLimitConfig `toml:"RateLimit"`
}

type CommonConfig struct {
	Port    int  `toml:"port"`
	OpenTLS bool `toml:"open_tls"`
}

func (config CommonConfig) Address() string {
	return fmt.Sprintf(":%d", config.Port)
}

type MySQLConfig struct {
	URL    string `toml:"url"`
	User   string `toml:"user"`
	Pwd    string `toml:"pwd"`
	DBName string `toml:"db_name"`
}

func (config MySQLConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Pwd,
		config.URL,
		config.DBName,
	)
}

type MongoDBConfig struct {
	URL      string `toml:"url"`
	User     string `toml:"user"`
	Pwd      string `toml:"pwd"`
	Database string `toml:"database"`
	MaxPool  uint64 `toml:"max_pool"`
	MinPool  uint64 `toml:"min_pool"`
}

type RedisConfig struct {
	URL string `toml:"url"`
	Pwd string `toml:"pwd"`
	DB  int    `toml:"db"`
}

type JWTConfig struct {
	Secret string `toml:"secret"`
	Expire int    `toml:"expire"`
	Enable bool   `toml:"enable"`
}

type RateLimitConfig struct {
	Enable        bool `toml:"enable"`
	WindowSeconds int  `toml:"window_seconds"`
	MaxRequests   int  `toml:"max_requests"`
}

func Init(configPath ...string) error {
	path := resolveConfigPath(firstPath(configPath))
	var loaded TomlConfig
	if _, err := toml.DecodeFile(path, &loaded); err != nil {
		return fmt.Errorf("decode config file %q: %w", path, err)
	}

	Conf = &loaded
	return nil
}

func firstPath(configPath []string) string {
	if len(configPath) == 0 {
		return ""
	}
	return configPath[0]
}

func resolveConfigPath(configPath string) string {
	if configPath != "" {
		if info, err := os.Stat(configPath); err == nil && info.IsDir() {
			return filepath.Join(configPath, configFileName())
		}
		return configPath
	}

	return filepath.Join("config", configFileName())
}

func configFileName() string {
	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = defaultEnvironment
	}
	return fmt.Sprintf("config-%s.toml", environment)
}

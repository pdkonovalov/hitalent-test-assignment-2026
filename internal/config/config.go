package config

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigEnv struct {
	LogLevel  string `env:"LOG_LEVEL" env-default:"info"`
	LogFormat string `env:"LOG_FORMAT" env-default:"text"`

	Endpoint string `env:"ENDPOINT" env-default:"localhost:8080"`
	Host     string `env:"HOST" env-default:"localhost"`
	Port     string `env:"PORT" env-default:"8080"`
	BasePath string `env:"BASE_PATH"`

	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresDB       string `env:"POSTGRES_DB"`
}

type Config struct {
	LogLevel  LogLevel
	LogFormat LogFormat

	Endpoint string
	Host     string
	Port     string
	BasePath string

	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	PostgresDB       string
}

func New() (*Config, error) {
	cfgEnv := ConfigEnv{}

	readConfigFileErr := cleanenv.ReadConfig(".env", &cfgEnv)
	if readConfigFileErr != nil {
		readEnvErr := cleanenv.ReadEnv(&cfgEnv)

		if readEnvErr != nil {
			return nil, errors.Join(readConfigFileErr, readEnvErr)
		}
	}

	cfg := Config{}

	logLevel, err := newLogLevel(cfgEnv.LogLevel)
	if err != nil {
		return nil, err
	}
	cfg.LogLevel = logLevel

	logFormat, err := newLogFormat(cfgEnv.LogFormat)
	if err != nil {
		return nil, err
	}
	cfg.LogFormat = logFormat

	cfg.Endpoint = cfgEnv.Endpoint
	cfg.Host = cfgEnv.Host
	cfg.Port = cfgEnv.Port
	cfg.BasePath = cfgEnv.BasePath

	cfg.PostgresUser = cfgEnv.PostgresUser
	cfg.PostgresPassword = cfgEnv.PostgresPassword
	cfg.PostgresHost = cfgEnv.PostgresHost
	cfg.PostgresPort = cfgEnv.PostgresPort
	cfg.PostgresDB = cfgEnv.PostgresDB

	return &cfg, nil
}

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelError LogLevel = "error"
)

func newLogLevel(s string) (LogLevel, error) {
	switch s {
	case "info":
		return LogLevelInfo, nil
	case "error":
		return LogLevelError, nil
	default:
		return "", fmt.Errorf("invalid log level: %s", s)
	}
}

func (l LogLevel) SlogLevel() slog.Level {
	switch l {
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormatJson LogFormat = "json"
)

func newLogFormat(s string) (LogFormat, error) {
	switch s {
	case "text":
		return LogFormatText, nil
	case "json":
		return LogFormatJson, nil
	default:
		return "", fmt.Errorf("invalid log format: %s", s)
	}
}

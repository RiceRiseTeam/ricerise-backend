package config

import (
	"ricerise/internal/logger"

	"github.com/go-playground/validator/v10"
	"github.com/samber/do/v2"
	"github.com/spf13/viper"
)

type AppConfig struct {
	AppPort int `mapstructure:"APP_PORT"`

	MySQLUrl      string `mapstructure:"MYSQL_URL"      validate:"required"`
	MySQLUser     string `mapstructure:"MYSQL_USER"     validate:"required"`
	MySQLPassword string `mapstructure:"MYSQL_PASSWORD" validate:"required"`

	RedisUrl      string `mapstructure:"REDIS_URL"      validate:"required"`
	RedisUser     string `mapstructure:"REDIS_USER"     validate:"required"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD" validate:"required"`

	AIBaseURL string `mapstructure:"AI_BASE_URL" validate:"required"`
	AIToken   string `mapstructure:"AI_TOKEN"    validate:"required"`

	JWTSecret string `mapstructure:"JWT_SECRET" validate:"required,min=32"`
	JWTExpire int    `mapstructure:"JWT_EXPIRE" validate:"required"`
}

func (c *AppConfig) validate() error {
	return validator.New().Struct(c)
}

func New(_ do.Injector) (*AppConfig, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		logger.Error("配置读取失败: %w", err)
		panic(err)
	}

	var config AppConfig
	if err := v.Unmarshal(&config); err != nil {
		logger.Error("配置解析失败: %w", err)
		panic(err)
	}

	if err := config.validate(); err != nil {
		logger.Error("非法配置项: %w", err)
		panic(err)
	}

	return &config, nil
}

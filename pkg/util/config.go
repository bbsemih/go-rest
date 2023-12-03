package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDR"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDR"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	RedisDB              int           `mapstructure:"REDIS_DB"`
	RedisPassword        string        `mapstructure:"REDIS_PASSWORD"`
	AWSRegion            string        `mapstructure:"AWS_REGION"`
	AWSBucketName        string        `mapstructure:"AWS_BUCKET_NAME"`
	RabbitMQURI          string        `mapstructure:"RABBITMQ_URI"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

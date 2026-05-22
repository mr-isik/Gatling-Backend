package infra

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Redis  RedisConfig
	Influx InfluxConfig
	LLM    LLMConfig
	JWT    JWTConfig
}

type JWTConfig struct {
	Secret           string
	AccessExpiryMin  int `mapstructure:"access_expiry_min"`
	RefreshExpiryDay int `mapstructure:"refresh_expiry_day"`
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	URL string
}

type InfluxConfig struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

type LLMConfig struct {
	APIKey string `mapstructure:"api_key"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

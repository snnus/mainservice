package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	Postgres    PgConfig    `yaml:"postgres"`
	Queueengine QeConfig    `yaml:"queueengine"`
	Kafka       KafkaConfig `yaml:"kafka"`
}

type KafkaConfig struct {
	Broker    string `yaml:"broker"`
	Topic     string `yaml:"topic"`
	BatchSize int    `yaml:"batch_size"`
}

type QeConfig struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

type PgConfig struct {
	Addr     string `yaml:"addr"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	DB       string `yaml:"db"`
	NShards  uint32 `yaml:"n_shards"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}

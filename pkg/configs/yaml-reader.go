package configs

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Http struct {
	Port int `yaml:"port"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
}

type Mongo struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type BucketPath struct {
	Raw       string `yaml:"raw"`
	Converted string `yaml:"converted"`
}

type YamlConfig struct {
	Http       Http       `yaml:"http"`
	Postgres   Postgres   `yaml:"postgres"`
	BucketPath BucketPath `yaml:"bucketPath"`
	MongoDB    Mongo      `yaml:"mongodb"`
}

func ReadYaml(path string) (*YamlConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config YamlConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

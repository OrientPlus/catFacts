package config

import (
	"caxfaxService1/internal/entity"
	"caxfaxService1/pkg/logger"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func LoadConfig(path string) (entity.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return entity.Config{}, fmt.Errorf("failed to open config file: %s", err.Error())
	}
	defer file.Close()

	var config entity.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return entity.Config{}, fmt.Errorf("failed to decode YAML: %s", err.Error())
	}

	config.Logger, err = logger.NewLogger(config.LoggerConfig)
	if err != nil {
		return entity.Config{}, fmt.Errorf("failed to initialize logger: %s", err.Error())
	}
	return config, nil
}

package entity

import "caxfaxService1/pkg/logger"

type Config struct {
	Logger             logger.Logger
	DbHost             string              `yaml:"db_host"`
	DbPort             string              `yaml:"db_port"`
	DbUser             string              `yaml:"db_user"`
	DbPassword         string              `yaml:"db_password"`
	DbName             string              `yaml:"db_name"`
	BrokerConsumerHost string              `yaml:"broker_consumer_host"`
	BrokerConsumerPort string              `yaml:"broker_consumer_port"`
	LoggerConfig       logger.LoggerConfig `yaml:"logger_config"`
}

package entity

import "caxfaxService2/pkg/logger"

type Config struct {
	Logger             logger.Logger
	RequestPeriod      int                 `yaml:"request_period_sec"`
	BrokerConsumerHost string              `yaml:"broker_consumer_host"`
	BrokerConsumerPort string              `yaml:"broker_consumer_port"`
	ResponseTopic      string              `yaml:"response_topic"`
	RequestTopic       string              `yaml:"request_topic"`
	LoggerConfig       logger.LoggerConfig `yaml:"logger_config"`
}

package service

import (
	"caxfaxService2/internal/broker"
	"caxfaxService2/internal/entity"
	"caxfaxService2/pkg/logger"
	"context"
	"fmt"
	"time"
)

type Service struct {
	requestPeriod time.Duration
	broker        *broker.Broker
	logger        logger.Logger
	stop          chan bool
}

// Выводит факт о котиках с заданным периодом
func (s *Service) Run() error {
	s.logger.Debug("Сервис запущен")
	for {
		select {
		case <-s.stop:
			s.logger.Info("Сервис остановлен")
			return nil

		default:
			time.Sleep(s.requestPeriod)

			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			fact, err := s.broker.RequestFact(ctx)
			if err != nil {
				s.logger.Errorf("Ошибка получения факта от сервиса: %s", err.Error())
				continue
			}

			fmt.Printf("Забавный факт:\n%s\n", fact.Message)
			s.logger.Infof("Получен факт %s; Длина: %d", fact.Message, fact.Length)
			continue
		}
	}
}

func (s *Service) Stop() {
	s.stop <- true
	s.broker.Stop()
}

func NewService(config entity.Config) (*Service, error) {
	service := &Service{}

	service.logger = config.Logger
	service.requestPeriod = time.Duration(config.RequestPeriod) * time.Second

	var err error
	service.broker, err = broker.NewBroker(config)
	if err != nil {
		return nil, fmt.Errorf("не удалось инициализировать брокер :%s", err.Error())
	}

	return service, nil
}

package service

import (
	"caxfaxService1/internal/api"
	"caxfaxService1/internal/broker"
	"caxfaxService1/internal/entity"
	"caxfaxService1/internal/repo"
	"caxfaxService1/pkg/logger"
	"context"
)

type Service1 struct {
	repo   repo.Repo
	api    api.APIClient
	broker *broker.Broker
	logger logger.Logger
	stopCh chan bool
}

func (s *Service1) Run() {
	s.broker.Run()
	s.logger.Debug("Сервис запущен")
	for {
		select {
		case request := <-s.broker.RequestCh:
			fact, err := s.api.GetFunFact()
			if err != nil {
				s.logger.Errorf("ошибка получения данных от внешнего API: %s", err.Error())
				continue
			}

			_, err = s.repo.AddFact(context.Background(), fact)
			if err != nil {
				s.logger.Errorf("ошибка добавления факта в базу: %s", err.Error())
				continue
			}

			err = s.broker.SendResponse(request, fact)
			if err != nil {
				s.logger.Errorf("ошибка отправки ответа в брокер сообщений: %s", err.Error())
			}
			continue

		case <-s.stopCh:
			s.logger.Debug("Сервис остановлен")
			return
		}
	}
}

func (s *Service1) Stop() {
	s.stopCh <- true
	s.broker.Stop()
	return
}

func NewService1(config entity.Config) (*Service1, error) {
	service := &Service1{}

	var err error
	service.repo, err = repo.NewRepo(config)
	if err != nil {
		return nil, err
	}

	service.api = api.NewCatFunFact(config.Logger)

	service.broker, err = broker.NewBroker(config)
	if err != nil {
		return nil, err
	}

	service.logger = config.Logger

	return service, nil
}

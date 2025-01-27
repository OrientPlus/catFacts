package main

import (
	"caxfaxService1/internal/service"
	"caxfaxService1/pkg/config"
	"fmt"
	"time"
)

func main() {
	time.Sleep(time.Second * 5)
	config, err := config.LoadConfig("./pkg/config/config.yml")
	if err != nil {
		fmt.Sprintf("Не удалось спарсить конфигурационный файл: %s", err.Error())
		return
	}

	service1, err := service.NewService1(config)
	if err != nil {
		fmt.Sprintf("Не удалось инициализировать сервис: %s", err.Error())
		return
	}

	service1.Run()
}

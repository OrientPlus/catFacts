package main

import (
	"caxfaxService2/internal/service"
	"caxfaxService2/pkg/config"
	"fmt"
	"time"
)

func main() {
	time.Sleep(time.Second * 7)
	config, err := config.LoadConfig("./pkg/config/config.yml")
	if err != nil {
		fmt.Printf("Не удалось спарсить конфигурационный файл: %s", err.Error())
		return
	}

	service, err := service.NewService(config)
	if err != nil {
		fmt.Printf("Не удалось запустить сервис: %s", err.Error())
		return
	}

	service.Run()

	return
}

# Сборка и запуск проекта

Для сборки и запуска проекта выполните следующие команды:

```bash
docker compose build
docker compose up
```

### Описание сервисов

## caxfaxService1
Ожидает запросы на получение факта через брокер сообщений.

После получения запроса:
* Обращается к внешнему API.
* Сохраняет полученный ответ в БД PostgreSQL.
* Отправляет ответ в caxfaxService2 через брокер сообщений Kafka.

## caxfaxService2
С заданным интервалом отправляет запросы сервису caxfaxService1 через брокера сообщений.
Получает ответ от caxfaxService1 и выводит его в stdout.

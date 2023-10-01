# Слушатель nats-streaming
### Описание
Сервис слушает канал nats-streaming. Полученные данные записывает в postgresDb

### Руководство по запуску
1. Поднять бд и nats (docker-compose файл)
2. Выполнить миграцию бд


    migrate -database postgres://postgres:password@localhost:5003/model?sslmode=disable -path db/migrations up
### Используемые технологии
    github.com/jmoiron/sqlx
    github.com/sirupsen/logrus
    github.com/spf13/viper
    github.com/nats-io/nats.go
    github.com/gin-gonic/gin
    github.com/swaggo/swag/cmd/swag
    
### Справочники
Документация nats:
https://docs.nats.io/running-a-nats-service/introduction/running
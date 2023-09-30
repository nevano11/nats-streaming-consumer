# Слушатель nats-streaming
### Описание
Сервис слушает канал nats-streaming. Полученные данные записывает в postgresDb

### Руководство по запуску
1. Поднять бд (docker-compose файл)
2. Выполнить миграцию


    migrate -database postgres://postgres:password@localhost:5003/model?sslmode=disable -path db/migrations up
3. натс??

### Используемые технологии
    github.com/jmoiron/sqlx
    github.com/sirupsen/logrus
    github.com/spf13/viper

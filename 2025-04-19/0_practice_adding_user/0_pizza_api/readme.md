# Практическое занятие Добавляем пользователей

## Системные требования
- golang 
- docker с docker compose (у меня также установлена новая версия сборщика образов `buildx`)
 
## Запуск проекта
### Для целей разработки переменные окружения хранятся в файле `.env` 
(а когда сервис разворачивается - уже через переменные окружения передаются настройки)
Для примера приложен файл `.env.example`, на его основе можно создать `.env`
```shell
cp .env.example .env
```

# Установка дополнительных инструментов (с помощью Makefile)
```shell
make install_tools
```

## Генерация спецификации из аннотаций к коду
```shell
swag init -g cmd/app/main.go -o api_docs -parseDependency
```

### Соберите из запустите проект в docker среде для разработки
```shell
source .env
make docker_alpine VERSION="${VERSION}"
```


# Запросы для тестов
### 1. Создание пользователя (POST /register)
```shell
source .env

curl -X POST -H "Content-Type: application/json" \
-d '{
  "name": "Иван Петров",
  "email": "ivan.petrov@mail.ru",
  "password": "SecureP@ss123",
  "phone": "+79001234567"
}' \
http://$HOST:$PORT/api/v1/register | jq
```

### Получение пользователей
```shell
source .env

curl -v http://$HOST:$PORT/api/v1/users/1 | jq
```

```shell
source .env

curl -v http://$HOST:$PORT/api/v1/users/99999 | jq
```

### Обновление пользователя
```shell
source .env

curl -vX PUT -H "Content-Type: application/json" \
 -d '{
    "name": "Иван Сергеевич Петров", 
    "phone": "+79009876543"
  }' \
  http://$HOST:$PORT/api/v1/users/1 | jq
```

### Удаление пользователя
```shell
source .env

curl -vX DELETE "http://$HOST:$PORT/api/v1/users/1"
```

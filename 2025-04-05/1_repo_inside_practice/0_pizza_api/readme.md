# Пример реализации Repository - разные механизмы внутренней реализации паттерна

## Запуск проекта
### Для целей разработки переменные окружения хранятся в файле `.env` 
(а когда сервис разворачивается - уже через переменные окружения передаются настройки)
Для примера приложен набор файлов `.env.example...`, на основе их можно создать `.env`
 ```shell
 cp .env.example.raw_sql .env
 ```

 ```shell
 cp .env.example.query_builder .env
 ```

 ```shell
 cp .env.example.orm .env
 ```

### Миграции
Установим инструмент для миграций
```shell
go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Перед запуском проекта - запустим миграции для создания первоначальной схемы базы данных
(только для типов репозитория `raw_sql` и `query_builder` )
```shell
source .env 
migrate -path=./migrations -database=sqlite://$DB_CONNECTION_STRING up
```
Убедитесь, что ошибок нет.

### Соберите и запустите проект
```shell
go build -o ./build/app ./cmd/app/main.go 
build/app
```

## Проверим с помощью запросов
### 1. Создание заказа (POST /orders)
```shell
source .env

curl -X POST -H "Content-Type: application/json" \
-d '{
"type": "margarita",
"size": "medium",
"quantity": 2,
"customer_id": 1
}' \
http://$HOST:$PORT/api/v1/orders | jq
```

```shell
source .env

curl -X POST -H "Content-Type: application/json" \
  -d '{
    "type": "pepperoni",
    "size": "large",
    "quantity": 1,
    "customer_id": 42
  }' \
  http://$HOST:$PORT/api/v1/orders | jq
```

### Получение заказов
```shell
source .env

curl http://$HOST:$PORT/api/v1/orders/1 | jq
```

```shell
source .env

curl -v http://$HOST:$PORT/api/v1/orders/99999 | jq
```

# Задание для самостоятельной работы

Добавьте недостающие методы - Update и Delete для репозитория, сервисов и контроллеров.
- При удалении данные не должны удаляться из БД, а лишь ставиться временная о удалении.
- При обновлении временная метка UpdatedAt дожна обновляться.
- Добавьте эти методы для всех трех типов репозитория.

*Добавьте метод GET с query параметрами для отбора и возможностью пагинации

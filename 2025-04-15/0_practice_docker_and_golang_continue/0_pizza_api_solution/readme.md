# Практическое занятие Docker и Golang

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

## Установим trivy
https://trivy.dev/latest/getting-started/installation/
> Всегда смотрите источник установки - особенно, когда запускаете скрипты из сети!
```shell
curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sudo sh -s -- -b /usr/local/bin v0.61.0
```


## Генерация спецификации из аннотаций к коду
```shell
swag init -g cmd/app/main.go -o api_docs -parseDependency
```
### Заполните переменные окружения (для примера можно взять файл `.env.examle`)
```shell
cp .env.example .env
```

## Проверим наш проект

### Проверим на соответствие лучшиим практикам создания `Dockerfile` с помощью `hadolint`:
```shell
docker run --rm -i hadolint/hadolint < Dockerfile
```

### Проверим на уязвимости с помощью trivy наш проект, передав путь к файловой системе
```shell
trivy fs .
```
Также он просканирует файл `go.mod` на версии пакетов с известными уязвимостями

### Соберите проект в docker
```shell
source .env
make docker_alpine VERSION="${VERSION}"
```

Проверка image с помощью trivy
```shell
source .env
trivy image "${IMAGE_NAME_ALPINE}:${VERSION}"
```


# Запросы для тестов
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

### Обновление заказа
```shell
source .env

curl -vX PUT -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "status": "delivering",
    "type": "pepperoni",
    "size": "large",
    "quantity": 2,
    "customer_id": 1
  }' \
  http://$HOST:$PORT/api/v1/orders/1 | jq
```

### Удаление заказа
```shell
source .env

curl -vX DELETE "http://$HOST:$PORT/api/v1/orders/1"
```

# Решение - что сделали
- `.env.example` - Поменяли ip адрес для того, чтобы можно было получить доступ к сервису в контейнере
- `Dockerfile*` (все Dockerfile'ы)
 - Используем `make` команду для того, чтобы заполнялись нужные нам флаги и переменные
 - Переформатировали команды для лучшей читаемости
 - Создаём папку `db/sqlite` - так как добавили её в файл `.dockerignore`, чтобы случайно не скопировать уже существующую БД
 - Изменили путь к бинарнику, так как сейчас компилируется с помощью команды `make`
- `.dockerignore` - для предотвращения копирования в образ ненужных файлов
- Проверили файлы `Dockerfile` с помощью `Hadolint` и `Trivy` на возможные ошибки и сделали необходимые изменения

# Решение - что сделали на втором занятии
- Исправили работу не под `root` - изменили права папкам. 
  - Проблему нашли с помощью инструмента `dive` - увидели, что владелец папки с базой данных - `root`
- Для контейнеров `scratch` и `distroless` разобрались как запускать не от `root`
  - Решение с копированием файлов с пользователями и группами (подходит и для других подобных ситуаций - пример `timezone` и `ca-certificates.crt`)
- Исправили `docker-compose.yml`
  - Пробросили порты наружу, чтобы можно было подключиться к сервису
  - Пробросили базу данных через `volume`
- Протестировали автоматическую перезагрузку приложения при изменении кода (инструмент `air`)
- Отладили код приложения в контейнере с помощью отладчика `delve`
- Посмотрели сборку контейнера для разработки с `air + delve` - hot reload + debugging
- Посмотрели разницу в слоях файловых систем контейнеров `scratch` и `debian`
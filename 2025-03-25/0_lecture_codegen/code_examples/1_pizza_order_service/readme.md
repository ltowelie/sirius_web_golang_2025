# Кодогенерация сервера с помощью oapi-codegen

## Установите кодогенератор
```shell
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

## Запустите генерацию кода для всех go файлов в текущей папке и всех подпапках 
(в этот раз используем Makefile для удобства)
```shell
make gen_api
```

## Убедитесь, что создался файл в папке `internal/web/api` `api.go`
```shell
ls -la internal/web/api | grep api.go
```

## Запустим наш сервер, но в начале пробросим переменные окружения
```shell
export HOST=localhost
export PORT=8080     
export LOGGER_LEVEL=D

go run ./cmd/server/main.go 
```

## В отдельном терминале можем протестировать запросы
### Проверка валидации параметров запроса
С правильными параметрами
```shell
curl "http://localhost:8080/orders?status=pending" 
```

С ошибкой в параметрах
```shell
curl "http://localhost:8080/orders?status=pendin" 
```
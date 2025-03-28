# Кодогенерация на примере генератора stringer

## Установите кодогенератор
```shell
go install golang.org/x/tools/cmd/stringer@latest
```

## Запустите генерацию кода для всех go файлов в текущей папке и всех подпапках
```shell
go generate ./...
```

`./...` - означает, что нужно искать и в подпапках

## Убедитесь, что создался файл в папке `internal/models/enums` `pizza_status_string.go`
```shell
ls -la internal/models/enums | grep pizza_status_string.go
```

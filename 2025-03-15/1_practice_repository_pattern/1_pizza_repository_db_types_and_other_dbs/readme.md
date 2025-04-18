# Пример реализации Repository на основе интерфейсов
Функционально реализовано создание слоя репозитория, кроме запросов к бд
Цель - показать студентам реализацию слоя репозитория на интерфейсах
Базовый функционал:
- Инициализация приложения
- Логгер
- Корректное завершение при поступлении сигнала - облегченная версия
- Сервис для демонстрации принципа работы с интерфейсами - передача репозитория

Для настройки используются переменные окружения:
1. `DB_TYPE` - выбор типа базы данных. В данный момент это sql или keyvalue
2. `DB_CONN_STR` - строка подключения к базе данных
3. `DB` - какую базу данных использовать (sqlite, postges при `DB_TYPE=sql`, leveldb - при `DB_TYPE=keyvalue`)

## Запуск проекта
1. Настройте переменные окружения
```shell
export DB=sqlite
export DB_CONN_STR=./sqlite
export DB_TYPE=sql
```
2. Соберите и запустите проект
```shell
go build -o ./build/app ./cmd/app/main.go 
./build/app
```

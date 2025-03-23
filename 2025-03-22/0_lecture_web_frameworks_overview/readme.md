---
date_created: "202503172313"
tags:
licence: "See LICENCE.md in root of this repository"
author: "Зайнутдинов Тимур Маратович"
---

# Лекция про различные подходы к написанию web сервисов, обзор фреймворков для написания web приложений

## Два способа написания web сервисов 
> на самом деле есть еще третий - но это уже за рамками нашей темы, чуть подробнее - в конце текста лекции

- Написание веб сервиса - логики и т.д. c использованием стандатрной библиотеки или фреймворков (gin, echo, fiber и т.д.). Во время написания сервиса - мы пишем специальные комментарии к функциям - хэндлерам. После - генерируем из комментариев файлы спецификации swagger.
- Можно назвать это подходом от обратного. Мы по спецификации swagger генерируем код сервера (можно и код клиента сгенерировать). Нам останется описать бизнес логику, а функции-хэндлеры будут автоматически сгенерированы.

Обычно второй подход используется реже.
На следующих практиках попробуем реализовать web сервис с использованием двух этих подходов.


## Обзор самых известных фрейморков для написания веб сервисов
Сейчас мы кратко пройдем по самым известным фреймворкам и не только фреймворкам

Но основная идея моя такая - принципиальной разницы нет в них. Когда вы устроитесь на работу бэкенд программистом - скорее всего в этой компании будет принято писать на одном из этих фреймворков. 
Хорошая новость в том, что итоговый код не очень сильно будет отличаться и вы быстро адаптируетесь к новому фрейворку.

Все примеры кода для фреймворков взяты из их репозиториев.

### Список фреймворков
- net/http
- Gin
- Echo
- Chi
- Fiber
- Gorilla Mux (не полноценный фрейворк, а скорее мультиплексер)
- Beego (много критики о нем)

### net/http
Мы уже писали несколько серверов с использованием http сервера из стандартной библиотеки.

На мой взгляд для простых сервисов - очень хороший выбор. 
Тем более, что большинство фрейворков указанных выше используют стандартный http-сервер. Просто у каких из проектов совместимость выше или ниже со стандартным пакетом http.

### Gin 
https://github.com/gin-gonic/gin
Использует `net/http`, но заменяет стандартный роутер на собственный (более быстрый). Middleware совместим с `net/http`.

Частично совместим с стандартным пакетом http.

Есть поддержка websocket, с использованием внешнего пакета
https://github.com/gin-gonic/examples/blob/master/websocket/server/server.go

Хэндлер GIN принимает контекст - `ctx`. В этом и есть несовместимость со стандартными функциями-хэндлерами.
Это аргумент, которое хранит всю информацию о запросе и имеет ряд методов, например, для записи ответов.
Мы у себя используем этот фреймворк. На дальнейших занятиях скорее всего буду писать код серверов на нём.

```go
router := gin.Default()
router.GET("/", func(c *gin.Context) { // GIN context
    c.String(200, "Hello")
})

// Запуск через стандартный http.Server
server := &http.Server{
    Handler: router,
    Addr:    ":8080",
}
server.ListenAndServe()
```

### Echo
https://github.com/labstack/echo

Полностью построен поверх `net/http`, сохраняет совместимость с его интерфейсами. Поддерживает стандартные middleware и обработчики.

Встроенная поддержка WebSocket
https://echo.labstack.com/docs/cookbook/websocket

Но как и в Gin - реализует свой собственный context (построен поверх стандартного `context.Context`) и функцию - хэндлер.

```go
package main

import (
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "log/slog"
  "net/http"
)

func main() {
  // Echo instance
  e := echo.New()

  // Middleware
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  // Routes
  e.GET("/", hello)

  // Start server
  // Под капотом запускает стандартный http сервер
  if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
    slog.Error("failed to start server", "error", err)
  }
}

// Handler
func hello(c echo.Context) error { // Функция - хэндлер не совместима со стандартной
  return c.String(http.StatusOK, "Hello, World!")
}
```

### Chi
https://github.com/go-chi/chi

100% совместимость с пакетом `net/http`. 
Работает как надстройка над `net/http`. Не заменяет, а расширяет его, предоставляя удобный роутер и middleware-подход.

```go
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// Функция-хэндлер полностью совместима с http пакетом
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":3000", r)
}
```

### Fiber
https://github.com/gofiber/fiber

Отличается от всех остальных проектов в данной подборке - он не совместим с стандартным http сервером. Под капотом использует https://github.com/valyala/fasthttp сервер, который по бенчмаркам быстрее стандартного сервера golang.
Не совместим с пакетом http
- Нет возможности использовать стандартные middleware
- Нет возможности использовать пакеты, совместимые с стандартным сервером, например go-swagger

Большой плюс - встроенная поддержка веб-сокетов

```go
app := fiber.New()

app.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("Hello")
})
app.Listen(":8080") // Использует fasthttp под капотом
```

### Gorilla Mux
https://github.com/gorilla/mux

Из названия понятно, что это не полноценный фреймворк, а мультиплексер. Он позволяет намного более удобно обрабатывать пути запросов.
В данную подборку попал потому, что приобрел широкую известность ранее. 
Если вы хотите полностью использовать только роутинг и собирать остальной стек самостоятельно.
Бывают проекты, где зависимости нужны минимальные, но возможностей встроенного мультиплексера в пакете http не хватает (хотя в последнее время разработчики go добавляют удобные штуки в стандартный http пакет)

Базовое использование не сильно отличается от стандартного мультиплексера. Но если изучить возможности - становится понятно, что мультиплексер очень удобный
```go
func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    r.HandleFunc("/products", ProductsHandler)
    r.HandleFunc("/articles", ArticlesHandler)
    http.Handle("/", r)
}
```

gorilla - это не один только мультиплексор, там целая группа инструментов для написания веб-сервисов
https://github.com/gorilla/


### beego
Но довольно популярен, если смотреть по звездам на github
Тянет на полноценный фреймворк (не goway). С ним ранее не сталкивался. 
По описанию функций - прям очень богатый. 
#### Функциональность (взято с github, ссылки сохранил)
[](https://github.com/beego/beego#features)
- RESTful support
- [MVC architecture](https://github.com/beego/beedoc/tree/master/en-US/mvc)
- Modularity
- [Auto API documents](https://github.com/beego/beedoc/blob/master/en-US/advantage/docs.md)
- [Annotation router](https://github.com/beego/beedoc/blob/master/en-US/mvc/controller/router.md)
- [Namespace](https://github.com/beego/beedoc/blob/master/en-US/mvc/controller/router.md#namespace)
- [Powerful development tools](https://github.com/beego/bee)
- Full stack for Web & API

#### Модули (взято с github, ссылки сохранил)
[](https://github.com/beego/beego#modules)

- [orm](https://github.com/beego/beedoc/tree/master/en-US/mvc/model)
- [session](https://github.com/beego/beedoc/blob/master/en-US/module/session.md)
- [logs](https://github.com/beego/beedoc/blob/master/en-US/module/logs.md)
- [config](https://github.com/beego/beedoc/blob/master/en-US/module/config.md)
- [cache](https://github.com/beego/beedoc/blob/master/en-US/module/cache.md)
- [context](https://github.com/beego/beedoc/blob/master/en-US/module/context.md)
- [admin](https://github.com/beego/beedoc/blob/master/en-US/module/admin.md)
- [httplib](https://github.com/beego/beedoc/blob/master/en-US/module/httplib.md)
- [task](https://github.com/beego/beedoc/blob/master/en-US/module/task.md)
- [i18n](https://github.com/beego/beedoc/blob/master/en-US/module/i18n.md)


## Про третий способ
Это не сказать, чтобы прям полноценный способ написать сервис, но иногда пригождается.
Если у нас есть `grpc` сервер в нашем сервисе и нам понадобилось его продублировать через `http` сервис - чтобы не писать заново всю логику, есть технология - `grpc gateway`. 
Она автоматически прокидывает сервис из `grpc` в `http` с помощью кодогенерации.

# Дополнительные материалы
- Статья про то как использовать grpc gateway вместо http фрейворков https://habr.com/ru/companies/ozonbank/articles/817381/
- Репозиторий с бенчмарками веб фрейворков на golang - https://github.com/smallnest/go-web-framework-benchmark
- Обзорная статья про beego - https://habr.com/ru/companies/otus/articles/802333/
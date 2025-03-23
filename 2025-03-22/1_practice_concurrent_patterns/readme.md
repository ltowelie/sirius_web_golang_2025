---
date_created: "202503160945"
tags: 
licence: See LICENCE.md in root of this repository
author: Зайнутдинов Тимур Маратович
---

# Практическое занятие по конкурентости (2 занятия)

Начнем с кода, представленного на докладе Роба Пайка о паттернах конкурентности

Вспомним материал, пройденный на лекции

Есть функция, которая долго и нудно выполняется, ожидание в рамках одной секунды
```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string) {
	for i := 0; ; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func main() {
	boring("test")
}
```

Запустим её в горутине, что произойдет с этой запущенной функцией?
```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string) {
	for i := 0; ; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func main() {
	go boring("test")
}
```

Добавим ожидание - пусть успеет отработать
```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string) {
	for i := 0; i < 4; i++ { // Поменять на 1000
		fmt.Println(msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func main() {
	go boring("test")
    // Тут подождем пока она закончится
    time.Sleep(5*time.Second)
	// Как считаете, хороший подход к ожиданию завершения горутины?
}
```

Чтобы вы изменили в этом коде, чтобы корректно дождаться завершения горутины (обсуждали на лекции механизмы)?

## Перейдем к простым паттернам

### Паттерн генератор (или когда удобно синхронное сообшение между горутинами)
Функции - объекты первого класса. В текущем примере функция будет возвращать канал.
А что это значит для нас, как программистов на golang? Их можно:
- Присваивать переменным (как в данном примере)
- Хранить в структурах данных
- Передавать в функции как аргументы
- Возвращать из функций
 
```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	c := boring("boring!") // Функция, которая возвращает канал.
	
	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c)
	}
	
	fmt.Println("You're boring; I'm leaving.")
}
func boring(msg string) <-chan string { // Возвращает канал только для чтения.
	c := make(chan string)
	
	go func() { // Запускаем горутину внутри функции.
		defer fmt.Println("Deferred exiting from goroutine")
		
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
		
		fmt.Println("Exiting from goroutine")
	}()
	
	return c // Возвращаем наш канал вызывающей стороне.
}
```

Все хорошо, код работает, но давайте рассмотрим его внимательнее. Давайте попробуем его сделать "production ready".

> [!question]- Какие проблемы в нем есть?

> [!answer]-  Ответ
> 1. Не закрываем канал
> 2. Не завершаем корректно горутину
> 3. Чтобы вы сюда добавили, чтобы выполнение программы работало корректно?


#### Первый способ - канал done (или Сигнальный канал)
Мы добавим в функцию второй канал на запись и вернем его из функции, назовём его `done`

- Способ с пометкой цикла
```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	c, done := boring("boring!") // Функция, которая возвращает канал.
	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c)
	}
	// Как лучше поступить?
	//done <- struct{}{} // Можно отправить сигнал
	close(done) // А можно просто закрыть
	time.Sleep(6*time.Second) // Ждём пока завершатся горутины
	fmt.Println("You're boring; I'm leaving.")
}

func boring(msg string) (<-chan string, chan<- struct{}) { // Возвращает канал только для чтения, и канал только для записи
	c := make(chan string)
	done := make(chan struct{}) // Канал для сигнализирования о завершении работы горутины

	go func() { 
		defer fmt.Println("Deferred exiting from goroutine")
		defer close(c) 
		// Пометка для нашего цикла, по ней мы можем прервать или перейти к следующей итерации
	loop: 
		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
				time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			case <-done: // Если приходит сигнал из канала, то завершаем цикл
				fmt.Println("Context is canceled")
				
				break loop // Останавливаем цикл с пометкой loop
				
				// Этот break сработает для select, а не для for
				// часто из-за этого возникают ошибки 
				// break 
			}
			
		}
		fmt.Println("Exiting from goroutine")
	}()

	return c, done // Добавили возврат канала `done`
}
```

- Способ с возвратом прямо из case
```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	c, done := boring("boring!") // Функция, которая возвращает канал.
	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c)
	}
	close(done)
	time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond) 
	fmt.Println("You're boring; I'm leaving.")
}

func boring(msg string) (<-chan string, chan<- struct{}) { // Возвращает канал только для чтения.
	c := make(chan string)
	done := make(chan struct{}) // Канал для сигнализирования о завершении работы горутины

	go func() { 
		defer fmt.Println("Deferred exiting from goroutine")
		defer close(c)
		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
				time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond) 
			case <-done:
				fmt.Println("Exiting from goroutine in case")
				
				return
			}
		}
		fmt.Println("Exiting from goroutine") // !!! Отличие от примера выше
	}()

	return c, done // Добавили возврат канала `done`
}
```

Используем пустую структуру `struct{}` для канала `done`, так как:
- Она не занимает память (`unsafe.Sizeof((struct{}) == 0`)
- Четко передает намерение: только сигнал, без данных

**Закрытие или Отправка значения в канал done**
Можно использовать и отправку значения в `done`, но 
- закрывать сигнальный канал все равно желательно. 
- придется отправлять на каждую горутину по сигналу.

## Второй способ - использование контекста

```go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background()) // Контекст с возвратом функции отмены
	defer cancel() // Важно вызвать для освобождения ресурсов

	c := boring(ctx, "boring!")
	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c)
	}
	cancel() // Сигнал остановки - вызываем второй раз
	
	fmt.Println("You're boring; I'm leaving.")
	time.Sleep(1000 * time.Millisecond)
}

func boring(ctx context.Context, msg string) <-chan string {
	c := make(chan string)

	go func() {
		defer close(c)

		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
				time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			
			// Как думаете, что возвращает метод Done()?
			case <-ctx.Done(): // Получаем сигнал отмены - через отмену контекста
				fmt.Println("boring: received cancel signal!")
			
				return
			}
		}
	}()

	return c
}

```

Что еще можно улучшить в данном коде?
>[!answer]- Ответ 
>Можно использовать WaitGroup для ожидания завершения горутины

## Мы можем пойти дальше и переиспользовать эту функцию
Можно представить, что эта функция представляет нам некий сервис.

Например, сервис поиска в интернете и через параметры мы можем задавать какой тип контента искать - видео, сайты, картинки.

А канал, который она возврашает - интерфейс для одностороннего (какого - получение или отправка из/в сервис?) взаимодействия с сервисом

Вот как это может выглядеть c двумя экземплярами "сервиса"
```go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем два "сервиса"
	// Здесь сервисы - это разные экземпляры в памяти или это один и тот же канал?
	pizzaMaker := boring(ctx, "Pizza is cooked and ready for delivery!")
	courier := boring(ctx, "I am delivery pizza!")
	
	
	for i := 0; i < 5; i++ {
		// Код в данном случае конкурентный?
		fmt.Printf("Pizza maker say: %q\n", <-pizzaMaker)
		fmt.Printf("Courier say: %q\n", <-courier)
	}
	
	cancel()
	
	fmt.Println("You're boring; I'm leaving.")
	time.Sleep(1000 * time.Millisecond)
}

func boring(ctx context.Context, msg string) <-chan string {
	c := make(chan string)

	go func() {
		defer close(c)

		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
				time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			case <-ctx.Done():
				fmt.Println("boring: received cancel signal!")
				
				return
			}
		}
	}()

	return c
}
```

Но тут мы вызываем сервисы последовательно и если первый сервис будет долго отвечать, то из второго сервиса мы прочитаем результат только после первого.
Было бы хорошо, если бы это было конкуренто причем результаты возвращались в главную горутину.

## Мультиплексирование каналов
Еще этот паттерн называется FanIn
```go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pizzaMaker := boring(ctx, "Pizza is cooked and ready for delivery!")
	courier := boring(ctx, "I am delivery pizza!")
	
    c := fanIn(pizzaMaker, courier) // Новая функция
    for i := 0; i < 5 ; i++ {
        fmt.Println(<-c)
    }
    
	cancel()
	
	fmt.Println("You're boring; I'm leaving.")
	time.Sleep(1000 * time.Millisecond)
}

func fanIn(input1, input2 <-chan string) <-chan string {
	// Создаем общий канал
	// В него будут конкуренто записываться результаты выполнения сервисов
    c := make(chan string)
    
    // Запускаем две горутины - по каждой на канал
    go func() { 
	    for { 
		    c <- <-input1 
		} 
	}()
        go func() { 
	    for { 
		    c <- <-input1 
		} 
	}()
	
    return c
}

func boring(ctx context.Context, msg string) <-chan string {
	c := make(chan string)

	go func() {
		defer close(c)

		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
				time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			case <-ctx.Done(): // Получаем сигнал отмены - через отмену контекста
				fmt.Println("boring: received cancel signal!")
				
				return
			}
		}
	}()

	return c
}
```

>[!question]- Можно ли как то улучшить данный паттерн?
Да, можно доработать так, чтобы он принимал неограниченное количество каналов.

>[!Answer]- Ответ
>```go
>
>func fanInMultichan(inputs ...<-chan string) <-chan string {
>    c := make(chan string) 
>    for _, input := range inputs {
>	    go func() { 
>		    for { 
>			    v, ok := <- input
>				if !ok {
>				
>					break
>				}
>			    c <- v 
>			} 
>		}()
 >   } 
>    	
 >   return c
>}
>```

## А еще можем улучшить наш FanIn с помощью select
```go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pizzaMaker := boring(ctx, "Pizza is cooked and ready for delivery!")
	courier := boring(ctx, "I am delivery pizza!")
	
	c := fanIn(ctx, pizzaMaker, courier)
	// Проверка выполнения при двух закрытых каналах
	// ставим пустой контекст, чтобы у нас горутина не завершилась 
	// после отмены контекста
    // c := fanIn(context.Background(), pizzaMaker, courier)
     
    for i := 0; i < 10; i++ {
        fmt.Println(<-c)
    }
    
	cancel() // Сигнал остановки
	
	fmt.Println("You're boring; I'm leaving.")
	time.Sleep(20000 * time.Millisecond)
}

func fanIn(ctx context.Context, input1, input2 <-chan string) <-chan string {
    c := make(chan string)
    go func() {
	    defer close(c) //
	    
        for {
            select {
            case s, ok := <-input1:  
	            if !ok { //
		            input1 = nil // 
		            
		            continue
	            }
	            c <- s
            case s, ok := <-input2:
              	if !ok { //
		            input2 = nil
		            
		            continue
	            }
	            c <- s
            case <-ctx.Done(): //
            
		        return
            }

			// Выполнится этот код, когда закроются каналы input1 и input2?
			if input1 == nil &&  input2 == nil { //
				fmt.Println("exit")
				
				return
			}

        }
    }()
    
    return c
}

func boring(ctx context.Context, msg string) <-chan string {
	c := make(chan string)

	go func() {
		defer close(c)

		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
				time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			case <-ctx.Done(): // Получаем сигнал отмены - через отмену контекста
				fmt.Println("boring: received cancel signal!")
				return
			}
		}
	}()

	return c
}
```

Главное преимущество данного подхода - благодаря `select` мы можем использовать всего одну горутину

>[!question] Можем ли мы доработать этот FanIn так, чтобы он принимал неограниченное количество каналов?

>[!answer] Ответ
>Так как мы динамически не можем создавать Case (если только через рефлексию - но это не очень хороший вариант) - то нет, сходу не приходит решения такой задачи: для заранее не заданного количества каналов создать `fanin` в одной горутине с использованием `select`

## Дополнительные материалы
- Статья про FanIn и его тестирование - https://habr.com/ru/articles/854302/
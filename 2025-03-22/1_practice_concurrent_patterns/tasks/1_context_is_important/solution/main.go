package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	orders := make(chan string)
	// Контекст для оповещения поваров о завершении работы
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Будем ждать поваров - чтобы они успели приготовить заказы
	// и чтобы не закрыть пиццерию раньше времени
	var wg sync.WaitGroup

	// Сегодня на кухне у нас будет трудиться три повара(горутины)
	// Запустите процесс готовки пицц - у нас сегодня как минимум 5 заказов
	for i := 1; i <= 3; i++ {
		name := fmt.Sprintf("Повар %d", i)
		wg.Add(1) // Увеличиваем счетчик горутин -
		go func(name string) {
			defer wg.Done() // Уменьшаем счетчик при завершении
			// Хорошая практика передавать контекст в функцию первым аргументом
			cook(ctx, name, orders)
		}(name)
	}

	// Отправляем заказы в планшеты поварам (в горутине)
	go func() {
		pizzas := []string{"Маргарита", "Пепперони", "Гавайская", "Четыре сыра", "Вегетарианская"}
		for _, pizza := range pizzas {
			orders <- pizza // Отправляем заказ в канал
		}
		close(orders) // Отправляем уведомление поварам - заказов больше нет
	}()

	// Даем пиццерии поработать 3 секунды перед закрытием
	time.Sleep(3 * time.Second)
	cancel() // Отправляем уведомление поварам - о завершении работы

	// Ждем пока все повара закончат работу
	wg.Wait()

	fmt.Println("Пиццерия закрыта!")
}

func cook(ctx context.Context, name string, orders <-chan string) {
	for {
		// Планшет поваров получает не только заказы,
		// но и информацию о завершении рабочего дня,
		// и информацию о том, что заказов больше нет.
		select {
		case order, ok := <-orders:
			if !ok {
				fmt.Printf("%s: заказов больше нет\n", name)

				return
			}
			fmt.Printf("%s начал готовить: %s\n", name, order)
			time.Sleep(3 * time.Second)
			fmt.Printf("%s закончил: %s\n", name, order)
		case <-ctx.Done():
			fmt.Printf("%s: получил сигнал закрытия пиццерии\n", name)

			return
		}
	}
}

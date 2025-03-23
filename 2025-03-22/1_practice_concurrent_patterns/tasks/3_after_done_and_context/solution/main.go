package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	orders := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		name := fmt.Sprintf("Повар %d", i)
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			cook(ctx, name, orders)
		}(name)
	}

	go func() {
		pizzas := []string{"Маргарита", "Пепперони", "Гавайская", "Четыре сыра", "Вегетарианская"}
		for _, pizza := range pizzas {
			orders <- pizza
		}
		close(orders)
	}()

	time.Sleep(1 * time.Second)
	cancel()

	wg.Wait()

	fmt.Println("Пиццерия закрыта!")
}

func cook(ctx context.Context, name string, orders <-chan string) {
	for {
		select {
		case order, ok := <-orders:
			if !ok {
				fmt.Printf("%s: заказов больше нет\n", name)

				return
			}
			fmt.Printf("%s начал готовить: %s\n", name, order)

			// Вложенный select, про который написал в ответе
			select {
			case <-time.After(3 * time.Second): // Сработает через три секунды
				fmt.Printf("%s закончил: %s\n", name, order)
			case <-ctx.Done(): // Не блокиреутся, если сигнал придет - завершим тут же работу
				fmt.Printf("%s получил сигнал закрытия пиццерии и отменил готовку пиццы: %s\n", name, order)

				return
			}

		case <-ctx.Done():
			fmt.Printf("%s: получил сигнал закрытия пиццерии\n", name)

			return
		}
	}
}

package main

import (
	"fmt"

	"stringer_enums/internal/models/enums"
)

func main() {
	// Запустим код до генерации и после.
	// Что будет выведено до генерации?
	// Как Println понимает, что ему нужно строку напечатать?
	fmt.Println(enums.StatusPreparing)
	fmt.Println(enums.StatusCancelled)

	status := enums.StatusDelivered
	fmt.Println(status.String()) // тут можно вызвать String()?
}

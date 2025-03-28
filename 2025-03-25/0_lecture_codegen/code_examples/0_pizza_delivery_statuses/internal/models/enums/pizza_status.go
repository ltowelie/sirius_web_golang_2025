package enums

// Ниже аннотация для команды go generate и параметры - какой генератор использовать (по идее просто каманда в shell)
//и, в данном случае, тип для генерации и срез префикса - чтобы в строках на было слова Status,
//и имя файла - куда выводить сгенерированный результат

//go:generate stringer -type=PizzaStatus -trimprefix=Status -output=pizza_status_string.go
type PizzaStatus int

const (
	StatusPending PizzaStatus = iota + 1 // Начинаем с 1
	StatusPreparing
	StatusInDelivery
	StatusDelivered
	StatusCancelled
)

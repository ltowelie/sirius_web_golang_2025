# Решение

Обратите внимание на
- Не всегда повар сразу завершает работу, а берет новый заказ
Как думаете почему так происходит?
- Логирование действий и сигналов дает нам полную картину происходящего
- Закрываем все каналы - так не будет утечки горутин
- Используем WaitGroup (или ErrGroup) для ожидания завершения работы поваров - не закроем пиццерию с поварами внутри
- Мы написали логику оповещения о завершении работы с использованием контекста с функцией отмены

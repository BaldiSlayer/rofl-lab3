# rofl-lab3. 

## Лабораторная работа 3. Фаззер

### Запуск

Синтаксис КС-грамматики
```
[rule] ::= [blank]*[NT][blank]*->[blank]*([NT][blank]*|[T][blank]*)+[EOL]+
[T] ::= [a - z]
[NT] ::= [A - Z][0 - 9]?|[[A - z]+ ([0 - 9])*]
```

```bash
go build ./cmd/main.go
./main --tests_count 10 --start_symbol S --break_prob 0.1 --terminal_adding_prob 0.2 < test
```

### Использованные материалы

[Приведение к НФХ](https://neerc.ifmo.ru/wiki/index.php?title=Нормальная_форма_Хомского)

https://neerc.ifmo.ru/wiki/index.php?title=Построение_FIRST_и_FOLLOW


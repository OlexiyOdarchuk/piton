# Importy

Імпорт у Piton робиться через `vykorystaty`.

## Import za ryadkom

```piton
vykorystaty "hello"
hello.Hello("svit")
```

## Import za identyfikatorom

```piton
vykorystaty hello
hello.Hello("svit")
```

Інтерпретатор доповнює шлях суфіксом `.piton`, виконує файл і створює модульне оточення.

## Yak tse pratsyuye

Після імпорту:

- модуль доступний під своїм ім'ям, наприклад `hello.Hello(...)`
- функції з модуля також можуть потрапити в глобальний простір, якщо там ще немає такого імені

## Pryklad iz repozytoriyu

Файл `examples/vykorystaty-demo.piton` імпортує локальний файл `hello.piton` із тієї ж директорії і звертається до функцій через модуль:

```piton
vykorystaty hello

functia Hi():
    drukuvaty "this is MY HI FUNCTIA"
    hello.Hello("world")
    hello.Hi("world!")

Hi()
```

## Obmezhennya

- шлях імпорту доповнюється `.piton`
- модульний механізм простий: це не пакетний менеджер і не ізольована система namespace-ів
- якщо файл не читається, інтерпретатор надрукує помилку читання
